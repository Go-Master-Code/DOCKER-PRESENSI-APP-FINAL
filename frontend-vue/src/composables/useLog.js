import {ref} from "vue";
import * as logService from "@/services/logService";

export function useLog() {
    const logList = ref([])
    const loading = ref(false);

    const fetchLogs = async () => {
        loading.value=true;
        try {
            const res = await logService.getAllLogs();
            logList.value = res.data.data || [] // jika res.data.data = null, maka return map kosong []
        } finally {
            loading.value = false;
        }
    }

    return { logList, loading, fetchLogs}
}