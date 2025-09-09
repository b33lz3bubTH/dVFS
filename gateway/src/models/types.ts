export interface StorageNode {
  id: string;
  baseUrl: string;
  healthEndpoint: string;
  isHealthy: boolean;
  lastChecked: Date;
}

export interface FileMetadata {
  id: string;
  filename: string;
  contentType: string;
  size: number;
  extension: string;
  nodeUrl: string;
  virtualPath: string;
  userId: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface Folder {
  id: string;
  name: string;
  path: string;
  userId: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface User {
  id: string;
  email: string;
  createdAt: Date;
  updatedAt: Date;
}