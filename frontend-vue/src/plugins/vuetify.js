import "vuetify/styles";
import { createVuetify } from "vuetify";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";

// 🔥 theme custom
const myCustomLightTheme = {
  dark: false,
  colors: {
    primary: '#1976D2',
    secondary: '#424242',
    background: '#F5F5F5',
  },
}

const myCustomDarkTheme = {
  dark: true,
  colors: {
    primary: '#90CAF9',
    secondary: '#BDBDBD',
    background: '#121212',
  },
}

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'light',

    themes: {
      light: myCustomLightTheme,
      dark: myCustomDarkTheme,
    },
  },
})

/* DEFAULT METHOD createVuetify
export default createVuetify({
  components,
  directives,
});
*/