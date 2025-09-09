import axios from 'axios';
import * as fs from 'fs';
import * as path from 'path';
import { StorageNode } from '../models/types';
import FormData from 'form-data';

export class NodePoolService {
  private nodes: StorageNode[] = [];
  private currentNodeIndex = 0;
  private readonly configPath: string;
  private readonly healthCheckInterval: number;
  private healthCheckPromise: Promise<void> | null = null;
  private isInitialized = false;

  constructor() {
    this.configPath = path.join(__dirname, '../../config/nodes.json');
    const config = JSON.parse(fs.readFileSync(this.configPath, 'utf-8'));
    this.healthCheckInterval = config.healthCheckIntervalMs || 10000;
  }

  async initialize(): Promise<void> {
    if (this.isInitialized) return;

    this.loadNodesFromConfig();
    await this.performHealthCheck(); // Initial health check
    this.startHealthCheck();
    this.isInitialized = true;
    
    const healthyNodes = this.getHealthyNodes();
    if (healthyNodes.length === 0) {
      console.warn('No healthy nodes available after initialization');
    }
  }

  private loadNodesFromConfig(): void {
    try {
      const config = JSON.parse(fs.readFileSync(this.configPath, 'utf-8'));
      this.nodes = config.nodes.map((node: any) => ({
        ...node,
        isHealthy: false,
        lastChecked: new Date()
      }));
      console.log(`Loaded ${this.nodes.length} nodes from config`);
    } catch (error: any) {
      console.error('Failed to load nodes from config:', error.message);
      this.nodes = [];
    }
  }

  private async checkNodeHealth(node: StorageNode): Promise<boolean> {
    try {
      const response = await axios.get(`${node.baseUrl}${node.healthEndpoint}`, {
        timeout: 5000 // 5 second timeout
      });
      return response.status === 200;
    } catch (error: any) {
      console.error(`Health check failed for node ${node.id}:`, error.message);
      return false;
    }
  }

  private async performHealthCheck(): Promise<void> {
    if (this.nodes.length === 0) {
      this.loadNodesFromConfig(); // Reload config if no nodes
    }

    console.log('Starting health check for all nodes...');
    const results = await Promise.all(
      this.nodes.map(async node => {
        const wasHealthy = node.isHealthy;
        const isHealthy = await this.checkNodeHealth(node);
        
        if (wasHealthy !== isHealthy) {
          console.log(`Node ${node.id} health status changed: ${isHealthy ? 'healthy' : 'unhealthy'}`);
        }

        node.isHealthy = isHealthy;
        node.lastChecked = new Date();
        return isHealthy;
      })
    );

    const healthyCount = results.filter(Boolean).length;
    console.log(`Health check complete. ${healthyCount}/${this.nodes.length} nodes are healthy`);
    this.healthCheckPromise = null;
  }

  private startHealthCheck(): void {
    setInterval(async () => {
      if (!this.healthCheckPromise) {
        this.healthCheckPromise = this.performHealthCheck();
        await this.healthCheckPromise.catch(err => {
          console.error('Health check failed:', err);
          this.healthCheckPromise = null;
        });
      }
    }, this.healthCheckInterval);
  }

  private getHealthyNodes(): StorageNode[] {
    const healthyNodes = this.nodes.filter(node => node.isHealthy);
    return healthyNodes;
  }

  private async selectHealthyNode(): Promise<StorageNode | null> {
    const healthyNodes = this.getHealthyNodes();
    if (healthyNodes.length === 0) {
      // Try immediate health check
      await this.performHealthCheck();
      const updatedHealthyNodes = this.getHealthyNodes();
      if (updatedHealthyNodes.length === 0) {
        return null;
      }
      return updatedHealthyNodes[0];
    }

    this.currentNodeIndex = (this.currentNodeIndex + 1) % healthyNodes.length;
    return healthyNodes[this.currentNodeIndex];
  }

  async uploadFile(file: Express.Multer.File): Promise<{ url: string, nodeUrl: string }> {
    if (!this.isInitialized) {
      await this.initialize();
    }

    const node = await this.selectHealthyNode();
    if (!node) {
      throw new Error('No healthy storage nodes available');
    }

    // Verify node health before upload
    const isHealthy = await this.checkNodeHealth(node);
    if (!isHealthy) {
      node.isHealthy = false;
      return this.uploadFile(file); // Retry with another node
    }

    const formData = new FormData();
    formData.append('file', file.buffer, file.originalname);

    try {
      const response = await axios.post(`${node.baseUrl}/api/v1/files`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
        timeout: 10000 // 10 second timeout
      });

      return {
        url: response.data.url,
        nodeUrl: node.baseUrl
      };
    } catch (error) {
      console.error(`Upload failed on node ${node.id}:`, error);
      node.isHealthy = false; // Mark node as unhealthy
      return this.uploadFile(file); // Retry with another node
    }
  }

  async deleteFile(fileId: string, nodeUrl: string): Promise<void> {
    try {
      await axios.delete(`${nodeUrl}/api/v1/files/${fileId}`, {
        timeout: 5000
      });
    } catch (error) {
      console.error(`Delete failed for file ${fileId}:`, error);
      throw error;
    }
  }
}