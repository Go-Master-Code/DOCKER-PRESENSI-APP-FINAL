import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import vuetify from "./plugins/vuetify";
import "@mdi/font/css/materialdesignicons.css";
import "@/styles/table.css" // import css global agar di tiap view tidak perlu tulis ulang kode css di <style scooped>

const app = createApp(App);

app.use(router);
app.use(vuetify);

app.mount("#app");