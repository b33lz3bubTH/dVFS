import { Router } from 'express';
import multer from 'multer';
import { FileController } from '../controllers/file.controller';

const router = Router();
const fileController = new FileController();
const upload = multer({ storage: multer.memoryStorage() });

// File routes
router.post('/files', upload.single('file'), fileController.uploadFile);
router.get('/files/:id', fileController.downloadFile);
router.get('/files/:id/info', fileController.getFileInfo);
router.head('/files/:id', fileController.checkFileExists);
router.delete('/files/:id', fileController.deleteFile);

export const fileRoutes = router;