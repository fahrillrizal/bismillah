// server/routes/index.js
const express = require('express');
const router = express.Router();
const { PrismaClient } = require('@prisma/client');
const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');

const prisma = new PrismaClient();

const verifyToken = (req, res, next) => {
  const authHeader = req.headers.authorization;
  
  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return res.status(401).json({ message: 'Unauthorized: No token provided' });
  }
  
  const token = authHeader.split(' ')[1];
  
  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET);
    req.user = decoded;
    next();
  } catch (error) {
    return res.status(401).json({ message: 'Unauthorized: Invalid token' });
  }
};

const isAdmin = (req, res, next) => {
  if (!req.user || !req.user.isAdmin) {
    return res.status(403).json({ message: 'Forbidden: Admin access required' });
  }
  
  next();
};


// POST /api/auth/login
router.post('/api/auth/login', async (req, res) => {
  try {
    const { username, password } = req.body;
    
    if (!username || !password) {
      return res.status(400).json({ message: 'Username and password are required' });
    }
    
    const user = await prisma.user.findUnique({
      where: { username }
    });
    
    if (!user) {
      return res.status(401).json({ message: 'Invalid credentials' });
    }
    
    const isMatch = await bcrypt.compare(password, user.password);
    if (!isMatch) {
      return res.status(401).json({ message: 'Invalid credentials' });
    }
    
    const token = jwt.sign(
      { id: user.id, username: user.username, isAdmin: user.isAdmin },
      process.env.JWT_SECRET,
      { expiresIn: process.env.JWT_EXPIRES_IN || '1d' }
    );
    
    res.json({
      token,
      user: {
        id: user.id,
        username: user.username,
        isAdmin: user.isAdmin
      }
    });
  } catch (error) {
    console.error('Login error:', error);
    res.status(500).json({ message: 'Server error' });
  }
});


// GET /api/links/grouped - Get all active links grouped by category
router.get('/api/links/grouped', async (req, res) => {
  try {
    const categories = await prisma.category.findMany({
      include: {
        links: {
          where: { isActive: true },
          orderBy: { order: 'asc' }
        }
      },
      orderBy: { order: 'asc' }
    });
    
    res.json(categories);
  } catch (error) {
    console.error('Error fetching grouped links:', error);
    res.status(500).json({ message: 'Error fetching links' });
  }
});

// GET /api/links - Get all links (admin)
router.get('/api/links', verifyToken, isAdmin, async (req, res) => {
  try {
    const links = await prisma.link.findMany({
      orderBy: [
        { categoryId: 'asc' },
        { order: 'asc' }
      ],
      include: {
        category: true
      }
    });
    
    res.json(links);
  } catch (error) {
    console.error('Error fetching links:', error);
    res.status(500).json({ message: 'Error fetching links' });
  }
});

// GET /api/categories - Get all categories (admin)
router.get('/api/categories', verifyToken, isAdmin, async (req, res) => {
  try {
    const categories = await prisma.category.findMany({
      orderBy: { order: 'asc' }
    });
    
    res.json(categories);
  } catch (error) {
    console.error('Error fetching categories:', error);
    res.status(500).json({ message: 'Error fetching categories' });
  }
});

// POST /api/categories - Create category (admin)
router.post('/api/categories', verifyToken, isAdmin, async (req, res) => {
  try {
    const { name, order } = req.body;
    
    if (!name) {
      return res.status(400).json({ message: 'Category name is required' });
    }
    
    const category = await prisma.category.create({
      data: {
        name,
        order: order || 0
      }
    });
    
    res.status(201).json(category);
  } catch (error) {
    console.error('Error creating category:', error);
    res.status(500).json({ message: 'Error creating category' });
  }
});

// PUT /api/categories/:id - Update category (admin)
router.put('/api/categories/:id', verifyToken, isAdmin, async (req, res) => {
  try {
    const { id } = req.params;
    const { name, order } = req.body;
    
    if (!name) {
      return res.status(400).json({ message: 'Category name is required' });
    }
    
    const category = await prisma.category.update({
      where: { id: parseInt(id) },
      data: {
        name,
        order: order || 0
      }
    });
    
    res.json(category);
  } catch (error) {
    console.error('Error updating category:', error);
    res.status(500).json({ message: 'Error updating category' });
  }
});

// DELETE /api/categories/:id - Delete category (admin)
router.delete('/api/categories/:id', verifyToken, isAdmin, async (req, res) => {
  try {
    const { id } = req.params;
    
    // Check if category has links
    const linksCount = await prisma.link.count({
      where: { categoryId: parseInt(id) }
    });
    
    if (linksCount > 0) {
      return res.status(400).json({ message: 'Cannot delete category with links' });
    }
    
    await prisma.category.delete({
      where: { id: parseInt(id) }
    });
    
    res.json({ message: 'Category deleted successfully' });
  } catch (error) {
    console.error('Error deleting category:', error);
    res.status(500).json({ message: 'Error deleting category' });
  }
});

// POST /api/links - Create link (admin)
router.post('/api/links', verifyToken, isAdmin, async (req, res) => {
  try {
    const { name, url, categoryId, order, isActive } = req.body;
    
    if (!name || !url || !categoryId) {
      return res.status(400).json({ message: 'Name, URL, and category are required' });
    }
    
    // Check if category exists
    const category = await prisma.category.findUnique({
      where: { id: parseInt(categoryId) }
    });
    
    if (!category) {
      return res.status(404).json({ message: 'Category not found' });
    }
    
    const link = await prisma.link.create({
      data: {
        name,
        url,
        categoryId: parseInt(categoryId),
        order: order || 0,
        isActive: isActive === true
      }
    });
    
    res.status(201).json(link);
  } catch (error) {
    console.error('Error creating link:', error);
    res.status(500).json({ message: 'Error creating link' });
  }
});

// PUT /api/links/:id - Update link (admin)
router.put('/api/links/:id', verifyToken, isAdmin, async (req, res) => {
  try {
    const { id } = req.params;
    const { name, url, categoryId, order, isActive } = req.body;
    
    if (!name || !url || !categoryId) {
      return res.status(400).json({ message: 'Name, URL, and category are required' });
    }
    
    // Check if category exists
    const category = await prisma.category.findUnique({
      where: { id: parseInt(categoryId) }
    });
    
    if (!category) {
      return res.status(404).json({ message: 'Category not found' });
    }
    
    const link = await prisma.link.update({
      where: { id: parseInt(id) },
      data: {
        name,
        url,
        categoryId: parseInt(categoryId),
        order: order || 0,
        isActive: isActive === true,
        updatedAt: new Date()
      }
    });
    
    res.json(link);
  } catch (error) {
    console.error('Error updating link:', error);
    res.status(500).json({ message: 'Error updating link' });
  }
});

// DELETE /api/links/:id - Delete link (admin)
router.delete('/api/links/:id', verifyToken, isAdmin, async (req, res) => {
  try {
    const { id } = req.params;
    
    await prisma.link.delete({
      where: { id: parseInt(id) }
    });
    
    res.json({ message: 'Link deleted successfully' });
  } catch (error) {
    console.error('Error deleting link:', error);
    res.status(500).json({ message: 'Error deleting link' });
  }
});

module.exports = router;