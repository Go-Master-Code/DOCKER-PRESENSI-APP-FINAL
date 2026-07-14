import api from "./api";

export const backupDatabase = async() => {
    return api.post(
        "/api/backup",
        {},
        {
            responseType: "blob",
        }
    )
}