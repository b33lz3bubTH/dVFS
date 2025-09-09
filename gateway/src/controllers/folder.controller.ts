import { Request, Response } from 'express';
import { PrismaClient } from '@prisma/client';
import path from 'path';

const prisma = new PrismaClient();

export class FolderController {
  async createFolder(req: Request, res: Response) {
    try {
      const { user } = req;
      const { path: folderPath } = req.body;

      if (!folderPath) {
        return res.status(400).json({ error: 'Folder path is required' });
      }

      // Validate absolute path
      if (!path.isAbsolute(folderPath)) {
        return res.status(400).json({ error: 'Folder path must be absolute' });
      }

      // Get or create user
      const userRecord = await prisma.user.upsert({
        where: { email: user.email },
        update: {},
        create: { email: user.email }
      });

      // Create folder
      const folder = await prisma.folder.create({
        data: {
          name: path.basename(folderPath),
          path: folderPath,
          userId: userRecord.id
        }
      });

      res.status(201).json(folder);
    } catch (error) {
      console.error('Create folder error:', error);
      res.status(500).json({ error: 'Failed to create folder' });
    }
  }

  async getTree(req: Request, res: Response) {
    try {
      const { user } = req;

      // Get or create user
      const userRecord = await prisma.user.upsert({
        where: { email: user.email },
        update: {},
        create: { email: user.email }
      });

      // Get all folders and files for user
      const [folders, files] = await Promise.all([
        prisma.folder.findMany({
          where: { userId: userRecord.id },
          orderBy: { path: 'asc' }
        }),
        prisma.file.findMany({
          where: { userId: userRecord.id },
          orderBy: { virtualPath: 'asc' }
        })
      ]);

      // Build tree structure
      const tree: any = {
        name: '/',
        type: 'folder',
        children: []
      };

      // Add folders to tree
      folders.forEach(folder => {
        const parts = folder.path.split('/').filter(Boolean);
        let current = tree;
        parts.forEach((part: any, index: number) => {
          let child = current.children.find((c: any) => c.name === part);
          if (!child) {
            child = {
              name: part,
              type: 'folder',
              children: []
            };
            current.children.push(child);
          }
          current = child;
        });
      });

      // Add files to tree
      files.forEach(file => {
        const parts = file.virtualPath.split('/').filter(Boolean);
        const fileName = parts.pop()!;
        let current = tree;
        parts.forEach(part => {
          let child = current.children.find((c: any) => c.name === part);
          if (!child) {
            child = {
              name: part,
              type: 'folder',
              children: []
            };
            current.children.push(child);
          }
          current = child;
        });
        current.children.push({
          name: fileName,
          type: 'file',
          id: file.id,
          size: file.size,
          contentType: file.contentType
        });
      });

      res.json(tree);
    } catch (error) {
      console.error('Get tree error:', error);
      res.status(500).json({ error: 'Failed to get folder tree' });
    }
  }
}