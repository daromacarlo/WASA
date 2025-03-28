import { createRouter, createWebHashHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import LoginView from '../views/LoginView.vue';
import RegisterView from '../views/RegisterView.vue';
import UserChatsView from '../views/UserChatsView.vue'; 
import CreateGroupView from '../views/CreateGroupView.vue'; 
import SearchUserView from '../views/SearchUserView.vue';
import ChatView from '../views/ChatView.vue';
import ModifyGroup from '../views/ModifyGroup.vue';
import ModifyUser from '../views/ModifyUser.vue';
import InoltroView from '../views/InoltroView.vue';

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: HomeView,  // La pagina iniziale
    },
    {
      path: '/login',
      name: 'Login',
      component: LoginView,  // La pagina di login
    },
    {
      path: '/register',
      name: 'Register',
      component: RegisterView,  // La pagina di registrazione
    },
    {
      path: '/wasachat/:nickname/chats',  // Nuova route per le chat dell'utente
      name: 'UserChats',
      component: UserChatsView,  // Componente che gestisce la visualizzazione delle chat
      props: true,  // Passa i parametri della route come props al componente
    },
    {
      path: '/wasachat/:nickname/chats/creategroup',
      name: 'CreateGroup',
      component: CreateGroupView,  
    },

    {
      path: '/wasachat/:nickname/chats/searchuser',
      name: 'SearchUser',
      component: SearchUserView,  
    },


    {
      path: '/wasachat/:nickname/chats/:chat',
      name: 'ChatView',
      component: ChatView,  
    },

    {
      path: '/wasachat/:nickname/chats/:chat/settings',
      name: 'ModifyGroup',
      component: ModifyGroup,  
    },

    {
      path: '/wasachat/:nickname/settings',
      name: 'ModifyUser',
      component: ModifyUser,  
    },

    {
      path: '/wasachat/:nickname/:chat/inoltro/:message',
      name: 'InoltroView',
      component: InoltroView,  
    },
  ],
});

export default router;