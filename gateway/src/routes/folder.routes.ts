import { Router } from 'express';
import { FolderController } from '../controllers/folder.controller';

const router = Router();
const folderController = new FolderController();

// Folder routes
router.post('/folders', folderController.createFolder);
router.get('/tree', folderController.getTree);

export const folderRoutes = router;