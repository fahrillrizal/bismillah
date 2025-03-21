const { PrismaClient } = require("@prisma/client");
const prisma = new PrismaClient();

// Get all links grouped by category (public endpoint)
exports.getGroupedLinks = async (req, res) => {
  try {
    const includeEmpty = req.query.includeEmpty === "true";

    let categories = await prisma.category.findMany({
      include: {
        links: {
          where: { isActive: true },
          orderBy: { order: "asc" },
        },
      },
      orderBy: { order: "asc" },
    });

    if (!includeEmpty) {
      categories = categories.filter((category) => category.links.length > 0);
    }

    res.json(categories);
  } catch (error) {
    console.error("Error fetching grouped links:", error);
    res.status(500).json({ message: "Error fetching links" });
  }
};

// Get active links grouped by category
const getActiveLinks = async (req, res) => {
  try {
    const links = await prisma.link.findMany({
      where: { isActive: true },
      orderBy: [{ category: "asc" }, { createdAt: "desc" }],
    });

    const shopeeLinks = links.filter((link) => link.category === "SHOPEE");
    const blibliLinks = links.filter((link) => link.category === "BLIBLI");
    const lazadaLinks = links.filter((link) => link.category === "LAZADA");
    const tiktokLinks = links.filter((link) => link.category === "TIKTOK");
    const tokopediaLinks = links.filter((link) => link.category === "TOKOPEDIA");

    res.json({
      shopeeLinks,
      blibliLinks,
      lazadaLinks,
      tiktokLinks,
      tokopediaLinks,
    });
  } catch (error) {
    console.error("Error fetching active links:", error);
    res.status(500).json({ message: "Error fetching links" });
  }
};

// Get link by ID
const getLinkById = async (req, res) => {
  try {
    const { id } = req.params;

    const link = await prisma.link.findUnique({
      where: { id: parseInt(id) },
    });

    if (!link) {
      return res.status(404).json({ message: "Link not found" });
    }

    res.json(link);
  } catch (error) {
    console.error("Error fetching link:", error);
    res.status(500).json({ message: "Error fetching link" });
  }
};

// Create new link
const createLink = async (req, res) => {
  try {
    const { name, url, category, isActive } = req.body;

    if (!name || !url || !category) {
      return res
        .status(400)
        .json({ message: "Name, URL, and category are required" });
    }

    const link = await prisma.link.create({
      data: {
        name,
        url,        
        category,
        isActive: isActive === true,
      },
    });

    res.status(201).json(link);
  } catch (error) {
    console.error("Error creating link:", error);
    res.status(500).json({ message: "Error creating link" });
  }
};

const updateLink = async (req, res) => {
  try {
    const { id } = req.params;
    const { name, url, category, isActive } = req.body;

    if (!name || !url || !category) {
      return res
        .status(400)
        .json({ message: "Name, URL, and category are required" });
    }

    const existingLink = await prisma.link.findUnique({
      where: { id: parseInt(id) },
    });

    if (!existingLink) {
      return res.status(404).json({ message: "Link not found" });
    }

    const link = await prisma.link.update({
      where: { id: parseInt(id) },
      data: {
        name,
        url,        
        category,
        isActive: isActive === true,
        updatedAt: new Date(),
      },
    });

    res.json(link);
  } catch (error) {
    console.error("Error updating link:", error);
    res.status(500).json({ message: "Error updating link" });
  }
};

// Delete link
const deleteLink = async (req, res) => {
  try {
    const { id } = req.params;

    const existingLink = await prisma.link.findUnique({
      where: { id: parseInt(id) },
    });

    if (!existingLink) {
      return res.status(404).json({ message: "Link not found" });
    }

    await prisma.link.delete({
      where: { id: parseInt(id) },
    });

    res.json({ message: "Link deleted successfully" });
  } catch (error) {
    console.error("Error deleting link:", error);
    res.status(500).json({ message: "Error deleting link" });
  }
};

module.exports = {
  getAllLinks,
  getActiveLinks,
  getLinkById,
  createLink,
  updateLink,
  deleteLink,
};
