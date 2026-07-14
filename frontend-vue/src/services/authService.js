import api from "./api";

export const login = (payload) => {
    return api.post("/api/login", payload)
}