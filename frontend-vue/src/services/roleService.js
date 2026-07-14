import api from "./api";

// GET all jenis role
export const getAllRole = () => api.get("/api/role");
