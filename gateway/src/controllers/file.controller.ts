import { Request, Response } from 'express';
import { PrismaClient } from '@prisma/client';
import { NodePoolService } from '../services/nodePool.service';
import path from 'path';

const prisma = new PrismaClient();
const nodePool = new NodePoolService();

export class FileController {
  async uploadFile(req: Request, res: Response) {
    try {
      const { user } = req;
      const file = req.file;
      const { virtualPath } = req.body;

      if (!file || !virtualPath) {
        return res.status(400).json({ error: 'File and virtualPath are required' });
      }

      // Validate absolute path
      if (!path.isAbsolute(virtualPath)) {
        return res.status(400).json({ error: 'Virtual path must be absolute' });
      }

      // Get or create user
      const userRecord = await prisma.user.upsert({
        where: { email: user.email },
        update: {},
        create: { email: user.email }
      });

      // Upload to storage node
      const { url, nodeUrl } = await nodePool.uploadFile(file);

      // Save metadata
      const fileRecord = await prisma.file.create({
        data: {
          filename: file.originalname,
          contentType: file.mimetype,
          size: file.size,
          extension: path.extname(file.originalname),
          nodeUrl: url,
          virtualPath,
          userId: userRecord.id
        }
      });

      res.status(201).json(fileRecord);
    } catch (error) {
      console.error('Upload error:', error);
      res.status(500).json({ error: 'Failed to upload file' });
    }
  }

  async downloadFile(req: Request, res: Response) {
    try {
      const { id } = req.params;
      const file = await prisma.file.findUnique({ where: { id } });

      if (!file) {
        return res.status(404).json({ error: 'File not found' });
      }

      // Redirect to storage node
      res.redirect(file.nodeUrl);
    } catch (error) {
      res.status(500).json({ error: 'Failed to download file' });
    }
  }

  async getFileInfo(req: Request, res: Response) {
    try {
      const { id } = req.params;
      const file = await prisma.file.findUnique({ where: { id } });

      if (!file) {
        return res.status(404).json({ error: 'File not found' });
      }

      res.json(file);
    } catch (error) {
      res.status(500).json({ error: 'Failed to get file info' });
    }
  }

  async checkFileExists(req: Request, res: Response) {
    try {
      const { id } = req.params;
      const file = await prisma.file.findUnique({ where: { id } });

      if (!file) {
        return res.status(404).json({ error: 'File not found' });
      }

      res.status(200).end();
    } catch (error) {
      res.status(500).json({ error: 'Failed to check file' });
    }
  }

  async deleteFile(req: Request, res: Response) {
    try {
      const { id } = req.params;
      const file = await prisma.file.findUnique({ where: { id } });

      if (!file) {
        return res.status(404).json({ error: 'File not found' });
      }

      // Delete from storage node
      await nodePool.deleteFile(id, file.nodeUrl);

      // Delete metadata
      await prisma.file.delete({ where: { id } });

      res.status(204).end();
    } catch (error) {
      res.status(500).json({ error: 'Failed to delete file' });
    }
  }
}