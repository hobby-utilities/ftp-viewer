import {createApp} from 'vue'
import App from './App.vue'
import router from "../router"
import {createPinia} from "pinia"

import PrimeVue from 'primevue/config'
import 'primevue/resources/themes/lara-light-blue/theme.css'
import 'primevue/resources/primevue.min.css'

const app = createApp(App)
const pinia = createPinia()

app.use(router)
app.use(pinia)
app.use(PrimeVue)

app.mount('#app')