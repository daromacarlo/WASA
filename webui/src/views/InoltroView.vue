<template>
  <div>
    <!-- Bottone Annulla -->
    <button @click="goBack" class="cancel-button">
      Annulla
    </button>

    <button @click="openNewChatModal" class="button">
      Nuova chat
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

    <!-- Modale nuova chat -->
    <div v-if="showAddMemberModal" class="modal">
      <div class="modal-content">
        <h3>Cerca una persona</h3>
        <input 
          v-model="newMemberName" 
          type="text" 
          placeholder="Inserisci il nickname dell'utente" 
          class="modal-input"
          @keyup.enter="forwardToNewUser"
          :disabled="addingMember"
        />
        <div class="modal-buttons">
          <button 
            @click="forwardToNewUser"
            class="modal-button confirm"
            :disabled="!newMemberName || addingMember"
          >
            <span v-if="addingMember">Inoltrando...</span>
            <span v-else>Inoltra</span>
          </button>
          <button @click="closeNewChatModal" class="modal-button cancel">Annulla</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      chats: [],
      loading: false,
      messageToForward: null,
      showAddMemberModal: false,
      newMemberName: "",
      addingMember: false
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
        console.error("Errore durante il caricamento delle conversazioni:", e);
        alert("Errore durante il caricamento delle conversazioni.");
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
      const messageId = this.$route.params.message;
      const destinationChatId = chat.chat_id;

      if (messageId) {
        try {
          this.loading = true;
          const response = await this.$axios.post(
            `/wasachat/${nickname}/inoltro/${destinationChatId}/messaggi/${messageId}`
          );

          if (response.status >= 200 && response.status < 300) {
            alert("Messaggio inoltrato con successo!");
            this.goBack();
          } else {
            console.error("Errore nella risposta del server:", response);
            alert("Errore durante l'inoltro del messaggio.");
          }
        } catch (error) {
          console.error("Errore durante l'inoltro del messaggio:", error);
          alert("Errore durante l'inoltro del messaggio.");
        } finally {
          this.loading = false;
        }
      } else {
        this.$router.push(`/wasachat/${nickname}/chats/${destinationChatId}`);
      }
    },

    openNewChatModal() {
      this.showAddMemberModal = true;
      this.newMemberName = "";
      this.addingMember = false;
    },

    closeNewChatModal() {
      this.showAddMemberModal = false;
      this.newMemberName = "";
      this.addingMember = false;
    },

    async forwardToNewUser() {
      if (!this.newMemberName) return;
      this.addingMember = true;
      const nickname = this.$route.params.nickname;
      const destinatario = this.newMemberName;
      const messaggio = this.$route.params.message;

      try {
        const response = await this.$axios.post(
          `/inoltro/${nickname}/a/${destinatario}/inoltro/messaggi/${messaggio}`
        );

        if (response.status >= 200 && response.status < 300) {
          alert("Messaggio inoltrato con successo!");
          this.closeNewChatModal();
          this.goBack();
        } else {
          console.error("Errore nella risposta del server:", response);
          alert("Errore durante l'inoltro del messaggio.");
        }
      } catch (error) {
        console.error("Errore durante l'inoltro del messaggio a nuovo utente:", error);
        alert("Errore durante l'inoltro del messaggio.");
      } finally {
        this.addingMember = false;
      }
    }
  }
};
</script>

<style scoped>
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

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 10px;
  width: 300px;
}

.modal-input {
  width: 100%;
  padding: 10px;
  margin: 10px 0;
}

.modal-buttons {
  display: flex;
  justify-content: space-between;
}

.modal-button {
  padding: 8px 15px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.modal-button.confirm {
  background-color: #4caf50;
  color: white;
}

.modal-button.cancel {
  background-color: #f44336;
  color: white;
}
</style>
