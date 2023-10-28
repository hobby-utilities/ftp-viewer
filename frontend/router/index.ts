import {createRouter, createWebHistory} from "vue-router";
import HelloWorld from "../src/components/HelloWorld.vue";

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'default-view',
            component: HelloWorld,
        }
    ]
})

export default router