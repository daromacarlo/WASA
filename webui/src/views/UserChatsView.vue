<template>
  <div>
    <!-- Bottone per creare un gruppo -->
    <button type="submit" :disabled="loading" @click="goToCreateGroup">
      {{ loading ? "Caricamento..." : "Crea Gruppo" }}
    </button>

    <!-- Bottone per cercare persone -->
    <button type="submit" :disabled="loading" @click="goToSearchUser">
      {{ loading ? "Caricamento..." : "Cerca Persone" }}
    </button>

    <button @click="goToModifyProfile" title="Modifica profilo">
          <span class="button-text">Modifica profilo</span>
    </button>

    <button @click="logout" title="Logout">
          <span class="logout">Logout</span>
    </button>

    <!-- Lista delle conversazioni -->
    <div class="chats-container">
      <h2>Le tue conversazioni:</h2>
      <ul v-if="chats.length > 0">
        <li v-for="chat in chats" :key="chat.chat_id" @click="goToChat(chat)">
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
      chats: [], // Lista delle conversazioni
      loading: false, // Stato di caricamento per i bottoni
      error: null, // Messaggio di errore
    };
  },
  async created() {
    await this.loadChats(); // Carica le conversazioni al momento della creazione del componente
  },
  methods: {
    // Carica le conversazioni
    async loadChats() {
      const nickname = this.$route.params.nickname;
      try {
        const response = await this.$axios.get(`/wasachat/${nickname}/chats`);
        console.log(response.data); // Stampa i dati ricevuti dal backend

        // Aggiungi l'intestazione Base64 se manca
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

    // Formatta il timestamp
    formatTime(time) {
      const date = new Date(time);
      return date.toLocaleString();
    },

    // Gestisce gli errori di caricamento delle immagini
    handleImageError(event) {
      console.error("Errore nel caricamento dell'immagine:", event);
      event.target.src = "https://via.placeholder.com/50"; // Immagine di fallback
    },

    goToModifyProfile() {
      const nickname = this.$route.params.nickname;
      this.$router.push(`/wasachat/${nickname}/settings`);
    },

    // Reindirizza alla pagina di creazione del gruppo
    goToCreateGroup() {
      const nickname = this.$route.params.nickname;
      this.$router.push(`/wasachat/${nickname}/chats/creategroup`);
    },

    // Reindirizza alla pagina di ricerca utenti
    goToSearchUser() {
      const nickname = this.$route.params.nickname;
      this.$router.push(`/wasachat/${nickname}/chats/searchuser`);
    },

    // Reindirizza alla vista corretta in base al tipo di chat
    goToChat(chat) {
      const nickname = this.$route.params.nickname;
        // Se è un gruppo, reindirizza a GroupView
        this.$router.push(`/wasachat/${nickname}/chats/${chat.chat_id}`);
    },

    logout(){
      this.$router.push(`/`)
    }
  },
};
</script>

<style scoped>
.chats-container {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
  text-align: left;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  padding: 10px;
  border-bottom: 1px solid #ccc;
  cursor: pointer; /* Cambia il cursore per indicare che l'elemento è cliccabile */
}

li:hover {
  background-color: #f5f5f5; /* Cambia il colore di sfondo al passaggio del mouse */
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

.logout{
  margin: 0;
  font-size: 0.9em;
  color: #e90b0b;
}

button {
  padding: 10px 20px;
  background-color: #7dac10;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  margin-bottom: 20px;
  margin-top: 20px;
  margin-left: 20px;
  margin-right: 20px;
}

button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.error-message {
  color: red;
  margin-top: 10px;
}
</style>