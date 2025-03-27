<template>
  <div>
    <!-- Bottone Annulla -->
    <button @click="goBack" class="cancel-button">
      Annulla
    </button>

    <!-- Lista delle conversazioni -->
    <div class="chats-container">
      <h2>Seleziona una chat per inoltrare il messaggio:</h2>
      <ul v-if="chats.length > 0">
        <li v-for="chat in chats" :key="chat.chat_id" @click="forwardToChat(chat)">
          <div class="chat-item">
            <img
              v-if="chat.foto"
              :src="chat.foto"
              class="chat-photo"
              @error="handleImageError"
            />
            <div v-else class="chat-photo-placeholder">Nessuna foto</div>
            <div class="chat-info">
              <p class="chat-name">{{ chat.nome }}</p>
              <p v-if="chat.ultimosnip" class="chat-last-message">{{ chat.ultimosnip }}</p>
              <p v-if="chat.time" class="chat-time">{{ formatTime(chat.time) }}</p>
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
      try {
        const response = await this.$axios.get(`/wasachat/${nickname}/chats`);
        this.chats = response.data.map(chat => {
          if (chat.foto && !chat.foto.startsWith('data:image')) {
            chat.foto = `data:image/jpeg;base64,${chat.foto}`;
          }
          return chat;
        });
      } catch (e) {
        this.error = "Errore durante il caricamento delle conversazioni.";
        console.error(e);
      } finally {
        this.loading = false;
      }
    },

    formatTime(time) {
      const date = new Date(time);
      return date.toLocaleString();
    },

    handleImageError(event) {
      console.error("Errore nel caricamento dell'immagine:", event);
      event.target.src = "https://via.placeholder.com/50";
    },

    goBack() {
      const { nickname, chat } = this.$route.params;
      this.$router.push(`/wasachat/${nickname}/chats/${chat}`);
    },

    async forwardToChat(chat) {
      const nickname = this.$route.params.nickname;
      const currentChatId = this.$route.params.chat;
      const messageId = this.$route.params.message;
      const destinationChatId = chat.chat_id;
      
      if (messageId) {
        try {
          this.loading = true;
          await this.$axios.post(
            `/wasachat/${nickname}/inoltra/${destinationChatId}/${messageId}`
          );
          
          // Ritorna automaticamente alla chat di origine dopo l'inoltro
          this.goBack();
          
          // Mostra un messaggio di successo (opzionale)
          this.$toast.success('Messaggio inoltrato con successo!');
        } catch (error) {
          console.error("Errore durante l'inoltro del messaggio:", error);
          this.error = "Errore durante l'inoltro del messaggio";
          this.$toast.error('Errore durante l\'inoltro del messaggio');
        } finally {
          this.loading = false;
        }
      } else {
        this.$router.push(`/wasachat/${nickname}/chats/${destinationChatId}`);
      }
    }
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