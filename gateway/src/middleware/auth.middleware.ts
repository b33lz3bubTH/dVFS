import { Request, Response, NextFunction } from 'express';

export const authMiddleware = (req: Request, res: Response, next: NextFunction) => {
  const userEmail = req.headers['user-email'];

  if (!userEmail || typeof userEmail !== 'string') {
    return res.status(401).json({ error: 'User email is required' });
  }

  // Add user to request
  req.user = { email: userEmail };
  next();
};