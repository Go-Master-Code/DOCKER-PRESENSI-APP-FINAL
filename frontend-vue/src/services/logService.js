import api from "./api";

// GET all hari logs
export const getAllLogs = () => api.get("/api/logs");