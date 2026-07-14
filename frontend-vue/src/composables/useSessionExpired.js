import { ref } from "vue";

// 🔥 state dialog
const showSessionExpired =
    ref(false);

export const useSessionExpired = () => {

    // tampilkan modal
    const openSessionExpired = () => {
        showSessionExpired.value = true;
    };

    // tutup modal
    const closeSessionExpired = () => {
        showSessionExpired.value = false;
    };

    return {
        showSessionExpired,
        openSessionExpired,
        closeSessionExpired,
    };
};