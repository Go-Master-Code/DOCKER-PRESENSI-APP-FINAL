import {ref} from "vue";
import * as userService from "@/services/userService";

export function useUser() {
    const userList = ref([])
    const loading = ref(false);

    const fetchUser = async () => {
        loading.value=true;
        try {
            const res = await userService.getAllUser();
            userList.value = res.data.data
        } finally {
            loading.value = false;
        }
    }

    return { userList, loading, fetchUser}
}