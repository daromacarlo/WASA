<template>
    <div class="group-messages-container">
      <div class="group-header">
        <h2 class="group-title">{{ groupName }}</h2>
        <p class="group-members">{{ membersCount }} membri</p>
      </div>
      
      <ul v-if="messages.length > 0" class="messages-list" ref="messagesList">
        <li
          v-for="message in messages"
          :key="message.message_id"
          class="message-item"
          :class="{ 'current-user-message': isCurrentUser(message.sender) }"
        >
          <div class="message-sender" v-if="!isCurrentUser(message.sender)">
            {{ message.sender }}
          </div>
          
          <div v-if="message.foto" class="message-photo-container">
            <img :src="message.foto" class="message-photo" @error="handleImageError" />
            <div class="message-meta">
              <p class="message-time">{{ formatTime(message.time) }}</p>
              <p v-if="isCurrentUser(message.sender)" class="message-status">
                {{ message.letto ? "✔️✔️" : "✔️" }}
              </p>
            </div>
          </div>
  
          <div v-else class="message-text-container">
            <p class="message-text">{{ message.text || "Nessun testo disponibile" }}</p>
            <div class="message-meta">
              <p class="message-time">{{ formatTime(message.time) }}</p>
              <p v-if="isCurrentUser(message.sender)" class="message-status">
                {{ message.letto ? "✔️✔️" : "✔️" }}
              </p>
            </div>
          </div>
        </li>
      </ul>
  
      <p v-else class="no-messages">Nessun messaggio nel gruppo</p>
  
      <div class="message-input-container">
        <input
          v-model="newMessage"
          type="text"
          placeholder="Scrivi un messaggio..."
          class="message-input"
          @keyup.enter="sendGroupMessage"
        />
        <button @click="selectPhoto" class="photo-button">📷</button>
        <button @click="sendGroupMessage" class="send-button">Invia</button>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        messages: [],
        loading: false,
        error: null,
        newMessage: "",
        newPhoto: null,
        currentUser: "",
        groupName: "Nome del Gruppo",
        membersCount: 0,
        groupMembers: []
      };
    },
    async created() {
      this.currentUser = this.$route.params.userId; // O il modo in cui ottieni l'utente corrente
      await this.loadGroupData();
      await this.loadGroupMessages();
    },
    methods: {
      async loadGroupData() {
        try {
          const groupId = this.$route.params.groupId;
          const response = await this.$axios.get(`/groups/${groupId}`);
          this.groupName = response.data.name;
          this.groupMembers = response.data.members;
          this.membersCount = this.groupMembers.length;
        } catch (error) {
          console.error("Errore nel caricamento dei dati del gruppo:", error);
        }
      },
      
      async loadGroupMessages() {
        const groupId = this.$route.params.groupId;
        try {
          this.loading = true;
          const response = await this.$axios.get(`/groups/${groupId}/messages`);
          
          this.messages = response.data.map(message => {
            if (message.foto && !message.foto.startsWith("data:image")) {
              message.foto = `data:image/jpeg;base64,${message.foto}`;
            }
            return message;
          });
  
          this.$nextTick(() => {
            this.scrollToBottom();
          });
        } catch (error) {
          console.error("Errore nel caricamento dei messaggi:", error);
        } finally {
          this.loading = false;
        }
      },
      
      async sendGroupMessage() {
        if (this.newMessage.trim() || this.newPhoto) {
          const groupId = this.$route.params.groupId;
          const messageData = {
            text: this.newMessage.trim(),
            foto: this.newPhoto || null,
            sender: this.currentUser
          };
  
          try {
            // Aggiunta ottimistica
            const newMessage = {
              message_id: Date.now(),
              sender: this.currentUser,
              text: this.newMessage.trim(),
              foto: this.newPhoto || null,
              time: new Date().toISOString(),
              letto: false,
              ricevuto: true
            };
            
            this.messages.push(newMessage);
            this.newMessage = "";
            this.newPhoto = null;
            
            // Invio al backend
            const response = await this.$axios.post(
              `/groups/${groupId}/messages`,
              messageData
            );
  
            // Aggiorna con i dati reali dal backend
            if (response.data.message_id) {
              newMessage.message_id = response.data.message_id;
            }
  
            this.$nextTick(() => {
              this.scrollToBottom();
            });
          } catch (error) {
            console.error("Errore nell'invio del messaggio:", error);
            this.messages = this.messages.filter(m => m.message_id !== newMessage.message_id);
          }
        }
      },
      
      selectPhoto() {
        const input = document.createElement("input");
        input.type = "file";
        input.accept = "image/*";
        input.onchange = async (event) => {
          const file = event.target.files[0];
          if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
              this.newPhoto = e.target.result;
              this.sendGroupMessage();
            };
            reader.readAsDataURL(file);
          }
        };
        input.click();
      },
      
      isCurrentUser(sender) {
        return sender === this.currentUser;
      },
      
      formatTime(time) {
        return new Date(time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
      },
      
      handleImageError(event) {
        event.target.src = "https://via.placeholder.com/150";
      },
      
      scrollToBottom() {
        const container = this.$refs.messagesList;
        if (container) container.scrollTop = container.scrollHeight;
      }
    }
  };
  </script>
  
  <style scoped>
  .group-messages-container {
    max-width: 800px;
    margin: 0 auto;
    height: 100vh;
    display: flex;
    flex-direction: column;
    background-color: #f5f5f5;
  }
  
  .group-header {
    padding: 15px;
    background-color: #fff;
    border-bottom: 1px solid #e0e0e0;
    text-align: center;
  }
  
  .group-title {
    margin: 0;
    font-size: 1.4rem;
    color: #333;
  }
  
  .group-members {
    margin: 5px 0 0;
    font-size: 0.9rem;
    color: #666;
  }
  
  .messages-list {
    flex: 1;
    padding: 15px;
    overflow-y: auto;
    list-style: none;
    margin: 0;
  }
  
  .message-item {
    margin-bottom: 15px;
    max-width: 70%;
  }
  
  .current-user-message {
    margin-left: auto;
  }
  
  .message-sender {
    font-size: 0.8rem;
    color: #555;
    margin-bottom: 3px;
    font-weight: bold;
  }
  
  .message-text-container, .message-photo-container {
    padding: 10px 15px;
    border-radius: 18px;
    position: relative;
  }
  
  .current-user-message .message-text-container,
  .current-user-message .message-photo-container {
    background-color: #dcf8c6;
  }
  
  .message-item:not(.current-user-message) .message-text-container,
  .message-item:not(.current-user-message) .message-photo-container {
    background-color: #fff;
  }
  
  .message-text {
    margin: 0;
    word-break: break-word;
  }
  
  .message-photo {
    max-width: 100%;
    border-radius: 10px;
    max-height: 300px;
  }
  
  .message-meta {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    margin-top: 5px;
  }
  
  .message-time {
    font-size: 0.7rem;
    color: #999;
    margin: 0 5px 0 0;
  }
  
  .message-status {
    font-size: 0.7rem;
    color: #999;
    margin: 0;
  }
  
  .no-messages {
    text-align: center;
    padding: 20px;
    color: #999;
  }
  
  .message-input-container {
    display: flex;
    padding: 10px;
    background-color: #fff;
    border-top: 1px solid #e0e0e0;
  }
  
  .message-input {
    flex: 1;
    padding: 10px 15px;
    border: 1px solid #e0e0e0;
    border-radius: 20px;
    margin-right: 10px;
  }
  
  .photo-button, .send-button {
    padding: 10px 15px;
    border: none;
    border-radius: 50%;
    background-color: #075e54;
    color: white;
    cursor: pointer;
  }
  
  .photo-button {
    margin-right: 10px;
  }
  </style>