<template>
    <div>
      <!-- Bottone Annulla -->
      <button @click="goBack" class="cancel-button">
        Indietro
      </button>
  
      <!-- Lista dei partecipanti -->
      <div class="chats-container">
        <h2>Partecipanti della chat:</h2>
        <ul v-if="chats.length > 0">
          <li v-for="chat in chats" :key="chat.chat_id" @click="forwardToChat(chat)">
            <div class="chat-item">
              <div class="chat-info">
                <p class="chat-name">{{ chat.Nickname}}</p>
              </div>
            </div>
          </li>
        </ul>
        <p v-else>Nessuna conversazione trovata.</p>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        chats: [],
        loading: false,
        error: null,
        messageToForward: null
      };
    },
    async created() {
      await this.loadChats();
      
      if (this.$route.params.message) {
        this.messageToForward = this.$route.params.message;
      }
    },
    methods: {
      async loadChats() {
        const nickname = this.$route.params.nickname;
        const gruppo = this.$route.params.chat;
        try {
          const response = await this.$axios.get(`/wasachat/${nickname}/utenti/gruppi/${gruppo}`); 
          this.chats = response.data.map(chat => {
            if (chat.foto && !chat.foto.startsWith('data:image')) {
              chat.foto = `data:image/jpeg;base64,${chat.foto}`;
            }
            return chat;
          });
        } catch (e) {
          this.error = "Errore durante il caricamento dei partecipanti.";
          console.error(e);
        } finally {
          this.loading = false;
        }
      },

      goBack() {
        const { nickname, chat } = this.$route.params;
        this.$router.push(`/wasachat/${nickname}/chats/${chat}`);
      },
    },
  };
  </script>
  
  <style scoped>
  /* Stili rimangono identici alla versione precedente */
  .chats-container {
    padding: 20px;
    max-width: 600px;
    margin: 0 auto;
    text-align: left;
  }
  
  .cancel-button {
    padding: 10px 20px;
    background-color: #f44336;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
    margin: 20px;
  }
  
  ul {
    list-style-type: none;
    padding: 0;
  }
  
  li {
    padding: 10px;
    border-bottom: 1px solid #ccc;
    cursor: pointer;
  }
  
  li:hover {
    background-color: #f5f5f5;
  }
  
  .chat-item {
    display: flex;
    align-items: center;
  }
  
  .chat-photo {
    width: 50px;
    height: 50px;
    border-radius: 50%;
    margin-right: 10px;
  }
  
  .chat-photo-placeholder {
    width: 50px;
    height: 50px;
    border-radius: 50%;
    background-color: #ccc;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 10px;
    color: #666;
    font-size: 0.8em;
  }
  
  .chat-info {
    flex-grow: 1;
  }
  
  .chat-name {
    font-weight: bold;
    margin: 0;
  }
  
  .chat-last-message {
    margin: 5px 0;
    color: #666;
  }
  
  .chat-time {
    margin: 0;
    font-size: 0.9em;
    color: #999;
  }
  
  .error-message {
    color: red;
    margin-top: 10px;
  }
  </style>