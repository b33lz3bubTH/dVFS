import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import { PrismaClient } from '@prisma/client';
import { fileRoutes } from './routes/file.routes';
import { folderRoutes } from './routes/folder.routes';
import { authMiddleware } from './middleware/auth.middleware';
import { errorHandler } from './middleware/error.middleware';
import { NodePoolService } from './services/nodePool.service';

const app = express();
export const prisma = new PrismaClient();
export const nodePool = new NodePoolService();

// Middleware
app.use(cors());
app.use(helmet());
app.use(express.json());
app.use(authMiddleware);

// Routes
app.use('/api/v1', fileRoutes);
app.use('/api/v1', folderRoutes);

// Error handling
app.use(errorHandler);

const PORT = process.env.PORT || 3000;

async function bootstrap() {
  try {
    await prisma.$connect();
    console.log('Database connected successfully');
    
    await nodePool.initialize();
    console.log('Node pool initialized');

    app.listen(PORT, () => {
      console.log(`Server is running on port ${PORT}`);
    });
  } catch (error) {
    console.error('Failed to start server:', error);
    process.exit(1);
  }
}

// Handle graceful shutdown
process.on('SIGTERM', async () => {
  console.log('SIGTERM signal received');
  await prisma.$disconnect();
  process.exit(0);
});

bootstrap();