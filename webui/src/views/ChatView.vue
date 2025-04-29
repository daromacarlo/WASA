<template>
  <div class="messages-container">
    <button @click="goHome()" class="group-action-button danger" title="Torna alla home">
      <span class="button-text">Menù principale</span>
    </button>
    
    <!-- Barra superiore per i gruppi -->
    <div v-if="isGroup" class="group-header">
      <div class="group-actions">
        <button @click="goToUpdateGroup" class="group-action-button" title="Modifica gruppo">
          <span class="button-text">Modifica</span>
        </button>
        <button @click="openAddMemberModal" class="group-action-button" title="Aggiungi persone">
          <span class="button-text">Aggiungi</span>
        </button>
        <button @click="showQuitModal = true" class="group-action-button danger" title="Esci dal gruppo">
          <span class="button-text">Esci</span>
        </button>
      </div>
    </div>

    <!-- Lista messaggi -->
    <ul v-if="messages.length > 0" class="messages-list" ref="messagesList">
      <li
        v-for="message in messages"
        :key="message.message_id"
        class="message-item"
        @click="openMessageModal(message)"
        :class="{ 
          'message-item-right': isCurrentUser(message.autore),
          'message-item-group': isGroup && !isCurrentUser(message.autore)
        }"
      >
        <!-- Mostra nome autore nei gruppi -->
        <div v-if="isGroup && !isCurrentUser(message.autore)" class="message-sender">
          {{ message.autore }}
        </div>

        <!-- Visualizzazione della risposta (se presente) -->
        <div v-if="message.risposta" class="message-reply-container">
          <div class="message-reply-preview">
            <span class="reply-label">Risposta a:</span>
            <!-- Mostra l'autore del messaggio originale nei gruppi -->
            <span v-if="isGroup" class="reply-author">{{ getOriginalMessageAuthor(message.risposta) }}</span>
            <div class="reply-content">
              {{ getOriginalMessageText(message.risposta) }}
            </div>
          </div>
        </div>

        <div v-if="message.inoltrato">
          <span class="inoltrato-label">inoltrato</span>
        </div>

        <!-- Messaggio di tipo foto -->
        <div v-if="message.foto" class="message-photo-container">
          <img
            :src="message.foto"
            class="message-photo"
            @error="handleImageError"
          />
          <div class="message-meta">
            <p v-if="message.time" class="message-time">{{ formatTime(message.time) }}</p>
            <p v-if="message.ricevuto && isCurrentUser(message.autore)" class="message-status">
              {{ message.letto ? "✔️✔️" : "✔️" }}
            </p>
          </div>
        </div>

        <!-- Messaggio di tipo testo -->
        <div v-else class="message-text-container">
          <p class="message-text">{{ message.text || "Nessun testo disponibile" }}</p>
          <div class="message-meta">
            <p v-if="message.time" class="message-time">{{ formatTime(message.time) }}</p>
            <p v-if="message.ricevuto && isCurrentUser(message.autore)" class="message-status">
              {{ message.letto ? "✔️✔️" : "✔️" }}
            </p>
          </div>
        </div>

        <!-- Visualizzazione reazioni -->
        <div v-if="message.commenti && message.commenti.length > 0" class="message-reactions">
          <div 
            v-for="(users, reaction) in groupReactionsByType(message.commenti)" 
            :key="reaction" 
            class="reaction-badge"
            :title="users.join(', ')"
          >
            <span class="reaction-emoji">{{ reaction }}</span>
            <span v-if="users.length > 1" class="reaction-count">{{ users.length }}</span>
          </div>
        </div>
      </li>
    </ul>
    <p v-else class="no-messages">Nessun messaggio trovato.</p>

    <!-- Barra input messaggi -->
    <div class="message-input-container">
      <input
        v-model="newMessage"
        type="text"
        placeholder="Scrivi un messaggio..."
        class="message-input"
        @keyup.enter="sendMessage"
      />
      <button @click="selectPhoto" class="photo-button">
        📷
      </button>
      <button @click="sendMessage" class="send-button">
        Invia
      </button>
    </div>

    <!-- Modale aggiunta membri -->
    <div v-if="showAddMemberModal" class="modal">
      <div class="modal-content">
        <h3>Aggiungi persone al gruppo</h3>
        <input 
          v-model="newMemberName" 
          type="text" 
          placeholder="Inserisci il nickname dell'utente" 
          class="modal-input"
          @keyup.enter="addUserToGroup(newMemberName)"
          :disabled="addingMember"
        />
        <div class="modal-buttons">
          <button 
            @click="addUserToGroup(newMemberName)" 
            class="modal-button confirm"
            :disabled="!newMemberName || addingMember"
          >
            <span v-if="addingMember">Aggiungendo...</span>
            <span v-else>Aggiungi</span>
          </button>
          <button @click="closeAddMemberModal" class="modal-button cancel">Annulla</button>
        </div>
      </div>
    </div>

    <div v-if="showQuitModal" class="modal">
      <div class="modal-content">
        <h3>Esci dal gruppo</h3>
        <p>Sei sicuro di voler uscire definitivamente da questo gruppo?</p>
        <div class="modal-buttons">
          <button @click="quitGroup" class="modal-button danger">Esci</button>
          <button @click="showQuitModal = false" class="modal-button cancel">Annulla</button>
        </div>
      </div>
    </div>

    <div v-if="showMessageModal" class="modal">
      <div class="modal-content">
        <h3>Azioni sul messaggio:</h3>
        <button @click="closeMessageModal" class="modal-button cancel">Annulla</button>
        <button @click="openanswereMessageModal(selectedMessage)" class="modal-button">Rispondi</button>
        <button @click="openCommentMessageModal(selectedMessage)" class="modal-button">Commenta</button>
        <button @click="goShermataInoltro(selectedMessage)" class="modal-button">Inoltra</button>
        <button 
          v-if="selectedMessage && isCurrentUser(selectedMessage.autore)" 
          @click="deleteMessage" 
          class="modal-button danger"
        >
          Elimina
        </button>
        <button 
          v-if="hasUserCommented(selectedMessage)"
          @click="deleteUserComment(selectedMessage)"
          class="modal-button danger"
        >
          Elimina mio commento
        </button>
      </div>
    </div>

    <div v-if="showanswereMessageModal" class="modal">
      <div class="modal-content-large">
        <h3>Rispondi al messaggio</h3>
        <button @click="closeanswereMessageModal" class="modal-button cancel">
          Annulla
        </button>
        <input
          v-model="ans"  
          type="text"
          placeholder="Scrivi una risposta..."
          class="message-input-ans"
          @keyup.enter="sendReplyMessage"  
        />
        <button @click="ansselectPhoto" class="photo-button">
          📷
        </button>
        <button @click="sendReplyMessage" class="send-button">Invia</button>
      </div>
    </div>

    <div v-if="showCommentMessageModal" class="modal">
      <div class="modal-content">
        <h3>Reagisci al messaggio</h3>
        <div class="reactions-grid">
          <button 
            v-for="reaction in reactions" 
            :key="reaction" 
            class="reaction-button"
            :class="{ 'active': hasUserReacted(selectedMessage, reaction) }"
            @click="toggleReaction(reaction)"
            :title="getReactionName(reaction)"
          >
            {{ reaction }}
          </button>
        </div>
        <button @click="closeCommentMessageModal" class="modal-button cancel">
          Annulla
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      reactions: ["❤️", "😂", "🥺", "👍", "😭"],
      messages: [],
      loading: false,
      error: null,
      newMessage: "",
      newPhoto: null,
      currentUser: "",
      selectedMessage: null,
      ans: "",
      ansphoto: "",
      
      // Dati gruppo
      isGroup: false,
      groupId: null,
      groupName: "",
      groupMembersCount: 0,
      
      // Modali
      showAddMemberModal: false,
      showEditModal: false,
      showQuitModal: false,
      showMessageModal: false,
      showanswereMessageModal: false,
      showCommentMessageModal: false,
      newMemberName: "",
      editedGroupName: "",
      addingMember: false,
    };
  },

  async created() {
    this.currentUser = this.$route.params.nickname;
    await this.checkIfGroup();
    await this.loadMessages();
  },

  methods: {
    // Metodi per le reazioni
    groupReactionsByType(comments) {
      const grouped = {};
      comments.forEach(comment => {
        if (!grouped[comment.reazione]) {
          grouped[comment.reazione] = [];
        }
        if (!grouped[comment.reazione].includes(comment.autore)) {
          grouped[comment.reazione].push(comment.autore);
        }
      });
      return grouped;
    },

    getReactionName(reaction) {
      const names = {
        "❤️":  "Mi piace",
        "😂": "Divertente",
        "🥺": "Carino",
        "👍": "OK",
        "😭": "Triste"
      };
      return names[reaction] || reaction;
    },

    hasUserReacted(message, reaction) {
      if (!message || !message.commenti) return false;
      return message.commenti.some(comment => 
        comment.autore === this.currentUser && comment.reazione === reaction
      );
    },

    hasUserCommented(message) {
      if (!message || !message.commenti) return false;
      return message.commenti.some(comment => comment.autore === this.currentUser);
    },

    async deleteUserComment(message) {
      if (!message || !this.hasUserCommented(message)) return;

      try {
        // Trova il commento dell'utente corrente
        const userComment = message.commenti.find(c => c.autore === this.currentUser);
        if (!userComment) return;

        // Chiamata API per eliminare il commento
        await this.$axios.delete(
          `/wasachat/${this.currentUser}/messaggi/${this.selectedMessage.message_id}`
        );

        // Aggiorna localmente la lista dei commenti
        const messageIndex = this.messages.findIndex(
          msg => msg.message_id === message.message_id
        );
        
        if (messageIndex !== -1) {
          this.messages[messageIndex].commenti = this.messages[messageIndex].commenti.filter(
            c => c.commento_id !== userComment.commento_id
          );
        }

        this.closeMessageModal();
      } catch (error) {
        console.error("Errore durante l'eliminazione del commento:", error);
        alert(error.response?.data?.message || "Si è verificato un errore durante l'eliminazione del commento");
      }
      this.$router.go()
    },

    async toggleReaction(reaction) {
      if (!this.selectedMessage) return;
      
      try {
        const chatId = this.$route.params.chat;
        const nickname = this.$route.params.nickname;
        const messageId = this.selectedMessage.message_id;
        
        const hasReacted = this.hasUserReacted(this.selectedMessage, reaction);
        
        await this.$axios.post(
          `/wasachat/${nickname}/messaggi/${messageId}`,
          { reazione: hasReacted ? null : reaction }
        );

        // Aggiorna localmente i commenti
        const messageIndex = this.messages.findIndex(
          msg => msg.message_id === messageId
        );
        
        if (messageIndex !== -1) {
          if (!this.messages[messageIndex].commenti) {
            this.$set(this.messages[messageIndex], 'commenti', []);
          }
          
          if (hasReacted) {
            // Rimuovi la reazione dell'utente
            this.messages[messageIndex].commenti = this.messages[messageIndex].commenti.filter(
              c => !(c.autore === this.currentUser && c.reazione === reaction)
            );
          } else {
            // Aggiungi la nuova reazione
            this.messages[messageIndex].commenti.push({
              autore: this.currentUser,
              reazione: reaction,
              commento_id: Date.now() // ID temporaneo, sarà sostituito dal backend
            });
          }
        }

        this.closeCommentMessageModal();
        
      } catch (error) {
        console.error("Errore durante l'invio della reazione:", error);
        alert(error.response?.data?.message || "Si è verificato un errore durante l'invio della reazione");
      }
    },

    // Metodi esistenti
    async loadMessages() {
      const chatId = this.$route.params.chat;
      try {
        this.loading = true;
        const response = await this.$axios.get(`/wasachat/${this.currentUser}/chats/${chatId}`);
        
        this.messages = response.data.map((message) => {
          if (message.foto && !message.foto.startsWith("data:image")) {
            message.foto = `data:image/jpeg;base64,${message.foto}`;
          }
          // Assicurati che commenti sia sempre un array
          message.commenti = message.commenti || [];
          return message;
        });

        this.$nextTick(() => {
          this.scrollToBottom();
        });
      } catch (e) {
        this.error = "Errore durante il caricamento dei messaggi.";
        console.error(e);
      } finally {
        this.loading = false;
      }
    },

    async deleteMessage() {
      if (!this.selectedMessage) return;
      
      const chatId = this.$route.params.chat;
      try {
        await this.$axios.delete(
          `/wasachat/${this.currentUser}/chats/${chatId}/messaggi/${this.selectedMessage.message_id}`
        );
        this.messages = this.messages.filter(
          msg => msg.message_id !== this.selectedMessage.message_id
        );
        this.closeMessageModal();
      } catch (error) {
        console.error("Errore:", error);
        alert("Si è verificato un errore durante l'eliminazione del messaggio");
      }
    },

    async sendMessage() {
      if (this.newMessage.trim() || this.newPhoto) {
        const messageData = {
          testo: this.newMessage.trim(),
          foto: this.newPhoto || "",
        };

        const chatId = this.$route.params.chat;

        try {
          const newMessage = {
            message_id: Date.now(),
            autore: this.currentUser,
            text: this.newMessage.trim(),
            foto: this.newPhoto || null,
            time: new Date().toISOString(),
            letto: false,
            ricevuto: false,
            commenti: []
          };
          
          this.messages.push(newMessage);
          this.newMessage = "";
          this.newPhoto = null;
          
          await this.$axios.post(
            `/wasachat/${this.currentUser}/chats/${chatId}`,
            messageData
          );

          this.$router.go();
          
        } catch (error) {
          console.error("Errore durante l'invio del messaggio:", error);
          this.messages = this.messages.filter(m => m.message_id !== newMessage.message_id);
          alert("Errore durante l'invio del messaggio. Riprova.");
        }
      }
    },

    openCommentMessageModal(message) {
      this.selectedMessage = message;
      this.showCommentMessageModal = true;
      this.showMessageModal = false;
    },

    closeCommentMessageModal() {
      this.showCommentMessageModal = false;
    },

    openanswereMessageModal() {
      if (!this.selectedMessage) return;
      this.showanswereMessageModal = true;
      this.showMessageModal = false;
    },

    closeanswereMessageModal() {
      this.showanswereMessageModal = false;
      this.ans = "";
      this.ansphoto = "";
    },

    async sendReplyMessage() {
      if (!this.selectedMessage) {
        alert("Seleziona un messaggio a cui rispondere");
        return;
      }

      if (!this.ans.trim() && !this.ansphoto) {
        alert("Inserisci un messaggio o seleziona una foto");
        return;
      }

      const chatId = this.$route.params.chat;
      const messageData = {
        testo: this.ans.trim(),
        foto: this.ansphoto || "",
      };

      try {
        const response = await this.$axios.post(
          `/wasachat/${this.currentUser}/risposta/chats/${chatId}/${this.selectedMessage.message_id}`,
          messageData
        );

        const newReply = {
          message_id: response.data.message_id,
          autore: this.currentUser,
          text: this.ans.trim(),
          foto: this.ansphoto || null,
          time: new Date().toISOString(),
          risposta: this.selectedMessage.message_id,
          letto: false,
          ricevuto: false,
          commenti: []
        };
        this.messages.push(newReply);

        this.$router.go();

      } catch (error) {
        console.error("Errore durante l'invio della risposta:", error);
        alert(error.response?.data?.message || "Errore durante l'invio della risposta");
      } finally {
        this.ans = "";
        this.ansphoto = "";
        this.closeanswereMessageModal();
      }
    },

    ansselectPhoto() {
      const input = document.createElement("input");
      input.type = "file";
      input.accept = "image/*";
      input.onchange = async (event) => {
        const file = event.target.files[0];
        if (file) {
          const reader = new FileReader();
          reader.onload = async (e) => {
            this.ansphoto = e.target.result;
            await this.sendReplyMessage();
          };
          reader.readAsDataURL(file);
        }
      };
      input.click();
    },

    selectPhoto() {
      const input = document.createElement("input");
      input.type = "file";
      input.accept = "image/*";
      input.onchange = async (event) => {
        const file = event.target.files[0];
        if (file) {
          const reader = new FileReader();
          reader.onload = async (e) => {
            this.newPhoto = e.target.result;
            await this.sendMessage();
          };
          reader.readAsDataURL(file);
        }
      };
      input.click();
    },

    formatTime(time) {
      const date = new Date(time);
      return date.toLocaleString();
    },

    handleImageError(event) {
      console.error("Errore nel caricamento dell'immagine:", event);
      event.target.src = "https://via.placeholder.com/150";
    },

    isCurrentUser(author) {
      return author === this.currentUser;
    },

    scrollToBottom() {
      const container = this.$refs.messagesList;
      if (container) {
        container.scrollTop = container.scrollHeight;
      }
    },

    async checkIfGroup() {
      const chatId = this.$route.params.chat;
      try {
        const response = await this.$axios.get(`/check/${chatId}`);
        this.isGroup = response.data.is_group;
        if (this.isGroup) {
          this.groupId = response.data.group_id;
        }
      } catch (error) {
        console.error("Errore nel verificare se è un gruppo:", error);
      }
    },

    async goToUpdateGroup() {
      this.$router.push({ 
          name: 'ModifyGroup', 
          params: { nickname: this.currentUser, chat: this.$route.params.chat } 
        });
    },

    async goShermataInoltro() {
      this.$router.push({ 
          name: 'InoltroView', 
          params: { nickname: this.currentUser, chat: this.$route.params.chat, message: this.selectedMessage.message_id } 
        });
    },

    async goHome() {
      this.$router.push({ 
          name: 'UserChats', 
          params: { nickname: this.currentUser} 
        });
    },

    async quitGroup() {
      const nickname = this.$route.params.nickname;
      const chat = this.$route.params.chat;
      try {
        await this.$axios.delete(`/wasachat/${nickname}/chats/${chat}`);
        this.$router.push({ 
          name: 'UserChats', 
          params: { nickname: this.currentUser } 
        });
      } catch (error) {
        console.error("Errore durante l'uscita dal gruppo:", error);
        alert("Si è verificato un errore durante l'uscita dal gruppo");
      } finally {
        this.showQuitModal = false;
      }
    },

    openAddMemberModal() {
      this.showAddMemberModal = true;
      this.newMemberName = "";
      this.addingMember = false;
    },

    closeAddMemberModal() {
      this.showAddMemberModal = false;
      this.newMemberName = "";
      this.addingMember = false;
    },

    async addUserToGroup(nickname) {
      if (!nickname || !nickname.trim()) {
        alert("Inserisci un nickname valido");
        return;
      }

      try {
        this.addingMember = true;
        const currentUser = this.$route.params.nickname;
        const chatId = this.$route.params.chat;
        
        await this.$axios.put(
          `/wasachat/${currentUser}/chats/gruppi/${chatId}/aggiungi`,
          { utente_da_aggiungere: nickname.trim() }
        );
        
        alert(`${nickname} è stato aggiunto al gruppo con successo!`);
        this.closeAddMemberModal();
        
      } catch (error) {
        console.error("Errore nell'aggiungere l'utente al gruppo:", error);
        this.handleAddMemberError(error);
      } finally {
        this.addingMember = false;
      }
    },

    handleAddMemberError(error) {
      let errorMessage = "Si è verificato un errore durante l'aggiunta al gruppo";
      if (error.response) {
        switch(error.response.status) {
          case 400: errorMessage = "Richiesta malformata"; break;
          case 401: errorMessage = "Non autorizzato"; break;
          case 403: errorMessage = "Non hai i permessi per aggiungere membri"; break;
          case 404: errorMessage = "Utente non trovato"; break;
          case 409: errorMessage = "L'utente è già nel gruppo"; break;
          case 500: errorMessage = "Errore lato server. Riprova più tardi"; break;
        }
      } else if (error.request) {
        errorMessage = "Impossibile connettersi al server. Verifica la tua connessione";
      }
      
      alert(errorMessage);
    },
    
    openMessageModal(message) {
      this.selectedMessage = message;
      this.showMessageModal = true;
    },
    
    closeMessageModal() {
      this.showMessageModal = false;
      this.selectedMessage = null;
    },

    getOriginalMessageText(replyId) {
      if (!replyId) return "(Nessun riferimento al messaggio)";
      
      const originalMessage = this.messages.find(msg => msg.message_id === replyId);
      
      if (!originalMessage) {
        return "(Messaggio eliminato)";
      }
      
      if (originalMessage.foto) {
        return "📷 [Immagine]";
      }
      
      return originalMessage.text 
        ? originalMessage.text.substring(0, 50) + (originalMessage.text.length > 50 ? "..." : "")
        : "(Messaggio senza testo)";
    },

    getOriginalMessageAuthor(replyId) {
      if (!this.isGroup || !replyId) return "";
      
      const originalMessage = this.messages.find(msg => msg.message_id === replyId);
      
      if (!originalMessage) {
        return "(Utente non più nel gruppo): ";
      }
      
      return this.isCurrentUser(originalMessage.autore) 
        ? "Te: " 
        : `${originalMessage.autore}: `;
    }
  }
};
</script>

<style scoped>
.messages-container {
  padding: 10px;
  max-width: 1500px;
  margin: 0 auto;
  background-color: #ffffff;
  border-radius: 1px;
  box-shadow: 0 4px 19px rgba(0, 0, 0, 0.1);
  padding-bottom: 2px;
  height: calc(100vh - 70px);
  overflow: hidden;
}

.messages-list {
  list-style-type: none;
  padding: 20px;
  margin: 10px;
  max-height: calc(100% - 120px);
  overflow-y: auto;
  scroll-behavior: smooth;
  width: 100%;
  box-sizing: border-box;
}

.message-item {
  display: flex;
  align-items: flex-start;
  padding: 16px;
  margin-bottom: 16px;
  background-color: #a1d1a1;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  max-width: 70%;
}

.message-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.message-item-right {
  margin-left: auto;
  background-color: #d1ecf1;
}

.message-item-group {
  margin-right: auto;
  background-color: #cadda5;
}

.message-sender {
  font-weight: bold;
  font-size: 10px;
  color: #41814c;
  margin-bottom: 50px; 
}

.message-photo-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.message-photo {
  max-width: 100%;
  height: auto;
  border-radius: 10px;
  margin-bottom: 10px;
  max-height: 300px;
}

.message-text-container {
  display: flex;
  flex-direction: column;
}

.message-text {
  margin: 10;
  color: #495057;
  font-size: 16px;
  margin-top: 5px;
  margin-left: 15px;
  line-height: 1.6;
}

.message-meta {
  display: flex;
  flex-direction: column;
  margin-top: 8px;
}

.message-time {
  margin: 0;
  color: #868e96;
  font-size: 12px;
}

.message-status {
  margin: 0;
  color: #868e96;
  font-size: 12px;
  font-style: italic;
}

.message-reactions {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 8px;
  padding-top: 4px;
  border-top: 1px solid rgba(0,0,0,0.05);
}

.reaction-badge {
  display: flex;
  align-items: center;
  background-color: #f0f2f5;
  border-radius: 10px;
  padding: 2px 6px;
  font-size: 12px;
  cursor: default;
  transition: background-color 0.2s;
}

.reaction-badge:hover {
  background-color: #e4e6eb;
}

.reaction-emoji {
  margin-right: 4px;
}

.reaction-count {
  font-size: 11px;
  color: #65676B;
}

.no-messages {
  text-align: center;
  color: #6c757d;
  font-size: 16px;
  padding: 24px;
  background-color: #f8f9fa;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.message-input-container {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #ffffff;
  border-top: 1px solid #e0e0e0;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
  max-width: 1500px;
  margin: 0 auto;
}

.message-input {
  flex-grow: 1;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  margin-right: 10px;
  font-size: 16px;
}

.photo-button,
.send-button {
  padding: 10px 15px;
  border: none;
  border-radius: 20px;
  background-color: #007bff;
  color: #ffffff;
  font-size: 16px;
  cursor: pointer;
  margin-left: 10px;
  margin-bottom: 20px; 
  margin-top: 20px;
}

.photo-button:hover,
.send-button:hover {
  background-color: #0056b3;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e0e0e0;
}

.group-actions {
  display: flex;
  gap: 10px;
}

.group-action-button {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 8px 12px;
  border: none;
  border-radius: 20px;
  background-color: #f0f0f0;
  color: #495057;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.group-action-button:hover {
  background-color: #e0e0e0;
}

.group-action-button.danger {
  color: #dc3545;
  background-color: #f8d7da;
}

.group-action-button.danger:hover {
  background-color: #f1b0b7;
}

.button-text {
  font-size: 14px;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 10px;
  padding: 20px;
  width: 90%;
  max-width: 400px;
  box-shadow: 0 4px 19px rgba(0, 0, 0, 0.1);
}

.modal-content-large {
  background-color: white;
  border-radius: 10px;
  padding: 30px;
  width: 100%;
  max-width:1000px;
  max-height: 1000px;
  box-shadow: 0 4px 19px rgba(0, 0, 0, 1);
}

.modal h3 {
  margin-top: 0;
  color: #495057;
}

.modal-input {
  width: 100%;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  margin: 10px 0;
  font-size: 16px;
}

.modal-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
  margin: 10px 0;
}

.modal-button {
  padding: 8px 15px;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
  margin: 10px 0;
}

.modal-button.confirm {
  background-color: #007bff;
  color: white;
}

.modal-button.cancel {
  background-color: #6c757d;
  color: white;
}

.modal-button.danger {
  background-color: #dc3545;
  color: white;
  margin-top: 10px;
}

.modal-button.confirm:hover {
  background-color: #0056b3;
}

.modal-button.cancel:hover {
  background-color: #5a6268;
}

.modal-button.danger:hover {
  background-color: #c82333;
}

.message-reply-container {
  border-left: 3px solid #4CAF50;
  padding-left: 8px;
  margin-bottom: 8px;
  color: #666;
  font-size: 0.9em;
}

.reply-label {
  font-weight: bold;
  margin-right: 5px;
}

.reply-author {
  color: #4CAF50;
}

.reply-content {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.reactions-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 10px;
  margin: 20px 0;
  justify-items: center;
}

.reaction-button {
  font-size: 24px;
  padding: 8px;
  border: none;
  background: none;
  cursor: pointer;
  transition: all 0.2s;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.reaction-button:hover {
  background-color: #f0f2f5;
  transform: scale(1.2);
}

.reaction-button.active {
  background-color: #e7f3ff;
  transform: scale(1.1);
}

.inoltrato-label {
  font-size: 0.8em;
  color: #666;
  font-style: italic;
  margin-bottom: 5px;
}

.message-input-ans {
  flex-grow: 1;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  margin-right: 10px;
  font-size: 16px;
  width: 100%;
  margin-bottom: 10px;
}
</style>