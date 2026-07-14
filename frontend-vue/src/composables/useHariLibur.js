import {ref} from "vue";
import * as hariLiburService from "@/services/hariLiburService";

export function useHariLibur() {
    const hariLiburList = ref([])
    const loading = ref(false);

    const fetchHariLibur = async () => {
        loading.value=true;
        try {
            const res = await hariLiburService.getAllHariLibur();
            hariLiburList.value = res.data.data
        } finally {
            loading.value = false;
        }
    }

    return { hariLiburList, loading, fetchHariLibur}
}