<template>
  <div :class="isGroup ? 'messages_container_group' : 'messages_container_private'">
    <button @click="goHome()" title="Go Back" class="goBack_btn">Go Back</button>
    <div v-if="isGroup" class="group_header">
      <div class="group_actions">
        <button @click="showParticipants" title="Show participants" class="btn">Show participants</button>
        <button @click="goToUpdateGroup" title="Modify group" class="btn">Modify group</button>
        <button @click="openAddMemberModal" title="Add user" class="btn">Add user</button>
        <button @click="openQuitModal" title="Quit group" class="quit_btn">Quit group</button>
      </div>
    </div>

    <ul v-if="messages.length > 0" class="messages_list" ref="messagesList">
      <li
        v-for="message in messages"
        :key="message.message_id"
        class="message_item"
        @click="openMessageModal(message)"
        :class="{ 
          'message_item_right': isCurrentUser(message.idautore),
          'message_item_group': isGroup && !isCurrentUser(message.idautore)
        }"
      >

        <div v-if="isGroup && !isCurrentUser(message.idautore)" class="message_sender">
          {{ message.autore }}
        </div>

        <div v-if="message.risposta" class="message_reply-container">
          <div class="message_reply-preview">
            <span class="reply-label">Answer to: </span>
            <span class="reply-author">{{ getOriginalMessageAuthor(message.risposta) }}</span>
            <div class="reply-content">
              {{ getOriginalMessageText(message.risposta) }}
            </div>
          </div>
        </div>

        <div v-if="message.inoltrato">
          <span class="forward_label">Forwarded</span>
        </div>

        <div v-if="message.foto" class="message_photo-container">
          <img
            :src="message.foto"
            class="message_photo"
            @error="handleImageError"
          />
          <div>
            <p v-if="message.time" class="message_time">{{ formatTime(message.time) }}</p>
            <p v-if="isCurrentUser(message.idautore)" class="message_status">
              {{ message.letto ? "read" : (message.ricevuto ? "received" : "sent") }}
            </p>
            <div v-if="message.commenti && message.commenti.length > 0" class="message_reactions">
              <div 
                v-for="(comment, index) in message.commenti" 
                :key="index" 
                class="reaction_badge_text"
                :title="comment.autore">
                <span class="reaction_emoji_text">{{ comment.reazione }}</span>
              </div>
            </div>
          </div>
        </div>

        <div v-else class="message_text-container">
          <p class="message_text">{{ message.text }}</p>
          <div>
            <p v-if="message.time" class="message_time">{{ formatTime(message.time) }}</p>
            <p v-if="isCurrentUser(message.idautore)" class="message_status">
              {{ message.letto ? "read" : (message.ricevuto ? "received" : "sent") }}
            </p>
            <div v-if="message.commenti && message.commenti.length > 0" class="message_reactions">
              <div 
                v-for="(comment, index) in message.commenti" 
                :key="index" 
                class="reaction_badge_photo"
                :title="comment.autore">
                <span class="reaction_emoji_photo">{{ comment.reazione }}</span>
              </div>
            </div>
          </div>
        </div>
      </li>
    </ul>
    
<ul v-if="messages.length > 0" class="messages_list" ref="messagesList">
</ul>
<p v-else class="no_messages">No messages yet.</p>

    <div :class="isGroup ? 'message_input_group' : 'message_input_private'">
      <input
        v-model="newMessage"
        type="text"
        placeholder="..."
        class="message_input"
        @keyup.enter="sendMessage"
      />
      <button @click="selectPhoto" class="btn_2">Send photo</button>
      <button @click="sendMessage" class="btn_2">Send</button>
    </div>

    <div v-if="showAddMemberModal" class="modal">
      <div class="modal_content">
        <h3>Add users to group</h3>
        <form @submit.prevent="addUserToGroup(newMemberName)">
          <input 
            v-model="newMemberName" 
            type="text" 
            placeholder="Insert a nickname" 
            class="modal_input"
            required
          />
          <button type="submit" class="modal_btn">Add</button>
          <button type="button" @click="closeAddMemberModal" class="modal_btn_gray">Go Back</button>
        </form>
      </div>
    </div>

    <div v-if="showQuitModal" class="modal">
      <div class="modal_content">
        <h3>Quit group</h3>
        <p>Are you sure?</p>
        <button @click="quitGroup" class="modal_btn_red">Quit</button>
        <button @click="showQuitModal = false" class="modal_btn_gray">Go Back</button>
      </div>
    </div>

    <div v-if="showMessageModal" class="modal">
      <div class="modal_content">
        <h3>Message actions:</h3>
        <button @click="closeMessageModal" class="modal_btn_gray">Go Back</button>
        <button @click="openanswereMessageModal(selectedMessage)" class="modal_btn">Answer</button>
        <button @click="openCommentMessageModal(selectedMessage)" class="modal_btn">Comment</button>
        <button @click="goToForwardView(selectedMessage)" class="modal_btn">Forward</button>
        <button 
          v-if="selectedMessage && isCurrentUser(selectedMessage.idautore)" 
          @click="deleteMessage" 
          class="modal_btn_red"
        >
          Delete
        </button>
        <button 
          v-if="hasUserCommented(selectedMessage)"
          @click="deleteUserComment(selectedMessage)"
          class="modal_btn_red"
        >
          Delete comment
        </button>
      </div>
    </div>

    <div v-if="showanswereMessageModal" class="modal">
      <div class="modal_content-large">
        <h3>Answer to message</h3>
        <button @click="closeanswereMessageModal" class="modal_btn_gray">Go Back</button>
        <input
          v-model="ans"  
          type="text"
          placeholder="..."
          class="message_input_ans"
          @keyup.enter="sendReplyMessage"  
        />
        <button @click="ansselectPhoto" class="modal_btn">Send photo</button>
        <button @click="sendReplyMessage" class="modal_btn">Send</button>
      </div>
    </div>

    <div v-if="showCommentMessageModal" class="modal">
      <div class="modal_content">
        <h3>Comment message</h3>
        <div class="reactions-grid">
          <button 
            v-for="reaction in reactions" 
            :key="reaction" 
            class="reaction_button"
            :class="{reaction : hasUserReacted(selectedMessage, reaction) }"
            @click="toggleReaction(reaction)"
          >
            {{ reaction }}
          </button>
        </div>
        <button @click="closeCommentMessageModal" class="modal_btn_gray">Go back</button>
      </div>
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
      selectedMessage: null,
      ans: "",
      ansphoto: "",
      isGroup: false,
      groupId: null,
      groupName: "",
      groupMembersCount: 0,
      showAddMemberModal: false,
      showEditModal: false,
      showQuitModal: false,
      showMessageModal: false,
      showanswereMessageModal: false,
      showCommentMessageModal: false,
      newMemberName: "",
      editedGroupName: "",
      chatId : this.$route.params.chat,
      currentUser : this.$route.params.nickname,
      reactions: [
                    "❤️",
                    "👽",
                    "😡",
                    "🫡",
                    "🤔",
                    "🫢",
                    "🔝",
                    "🇮🇹",
                    "🗿",
                    "🤮",
                    "😱",
                    "🤓",
                    "😂",
                    "🥺",
                    "👍",
                    "😭"
                  ],
    };
  },

  async created() {

  try {
    const response = await this.$axios.get(`/wasachat/${this.currentUser}`);
    this.currentUserId = response.data.id;
    await this.checkIfGroup();
    await this.loadMessages();
  } catch (error) {
    console.error("Error:", error);
  }
},

  methods: {

  goHome(){
    this.$router.push(`/wasachat/${this.currentUser}/chats`);
      },

  isCurrentUser(idauthor) {
    if (idauthor == this.currentUserId){
      return true
    }
    else{
      return false
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
    this.showanswereMessageModal = true;
    this.showMessageModal = false;
  },

  closeanswereMessageModal() {
    this.showanswereMessageModal = false;
    this.ans = "";
    this.ansphoto = "";
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

  openMessageModal(message) {
    this.selectedMessage = message;
    this.showMessageModal = true;
  },
  
  closeMessageModal() {
    this.showMessageModal = false;
    this.selectedMessage = null;
  },

  openQuitModal() {
    this.showQuitModal = true;
  },

  closeQuitModal() {
    this.showQuitModal = false;
  },

  getOriginalMessageText(replyId) {
  let originalMessage = null;
  for (let i = 0; i < this.messages.length; i++) {
    if (this.messages[i].message_id == replyId) {
      originalMessage = this.messages[i];
      break;
    }
  }
  if (!originalMessage) {
    return "(Erased message)";
  }
  if (originalMessage.foto) {
    return "img";
  }
  if (originalMessage.text) {
    const text = originalMessage.text;
    if (text.length > 15) {
      return text.substring(0, 15) + "...";
    }
    return text;
    }
  },

  getOriginalMessageAuthor(replyId) {
  let originalMessage = null;
  for (let i = 0; i < this.messages.length; i++) {
    if (this.messages[i].message_id == replyId) {
      originalMessage = this.messages[i];
      break;
    }
  }
    if (!originalMessage) {
      return "(Erased): ";
    }
    else{
      return originalMessage.autore
    } 
   },

   //from stackOverflow
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
   //from stackOverflow
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

  async checkIfGroup() {
    try {
      const response = await this.$axios.get(`/check/${this.chatId}`);
      this.isGroup = response.data.is_group;
      if (this.isGroup) {
        this.groupId = response.data.group_id;
      }
    } catch (error) {
      if (error.response) {
        const messaggio = error.response.data.errore;
        const codiceErrore = parseInt(error.response.data.codiceErrore);
        alert(messaggio + ` (codice ${codiceErrore})`);
      } else {
        alert("Error: Network error.");
      }
    }
  },

  hasUserReacted(message, reaction) {
  for (let i = 0; i < message.commenti.length; i++) {
    if (message.commenti[i].autore == this.currentUser && 
        message.commenti[i].reazione == reaction) {
      return true;
      }
   }
    return false;
  },

  hasUserCommented(message) {
    for (let i = 0; i < message.commenti.length; i++) {
      if (message.commenti[i].idautore == this.currentUserId) {
        return true;
      }
    }
    return false;
  },

  async deleteUserComment(message) {
  try {
    let userComment = null;
    for (let i = 0; i < message.commenti.length; i++) {
      if (message.commenti[i].idautore == this.currentUserId) {
        userComment = message.commenti[i];
        break;
      }
    }
    
    if (!userComment) return;

    const response = await this.$axios.delete(
      `/wasachat/${this.currentUser}/messaggi/${this.selectedMessage.message_id}`
    );

    let messageIndex = -1;
    for (let i = 0; i < this.messages.length; i++) {
      if (this.messages[i].message_id == message.message_id) {
        messageIndex = i;
        break;
      }
    }

    if (messageIndex !== -1) {
      const filteredComments = [];
      for (let i = 0; i < this.messages[messageIndex].commenti.length; i++) {
        if (this.messages[messageIndex].commenti[i].commento_id !== userComment.commento_id) {
          filteredComments.push(this.messages[messageIndex].commenti[i]);
        }
      }
      this.messages[messageIndex].commenti = filteredComments;
    }

    this.$router.go();
    this.closeMessageModal();
  } catch (error) {
    if (error.response) {
      const messaggio = error.response.data.errore;
      const codiceErrore = parseInt(error.response.data.codiceErrore);
      alert(messaggio + ` (codice ${codiceErrore})`);
    } else {
      alert("Error: Network error.");
    }
    this.$router.go();
   }
  },

  async goToUpdateGroup() {
  this.$router.push({ 
      name: 'ModifyGroup', 
      params: { nickname: this.currentUser, chat: this.$route.params.chat } 
    });
  },
  
  async toggleReaction(reaction) {
  try {
    const messageId = this.selectedMessage.message_id;
    const hasReacted = this.hasUserReacted(this.selectedMessage, reaction);

    if (hasReacted) {
      this.closeCommentMessageModal();
      return;
    }
    
    const response = await this.$axios.post(
      `/wasachat/${this.currentUser}/messaggi/${messageId}`,
      { reazione: hasReacted ? null : reaction }
    );

    for (let i = 0; i < this.messages.length; i++) {
      const msg = this.messages[i];

      if (msg.message_id == messageId) {
        if (!msg.commenti) {
          this.$set(this.messages[i], 'commenti', []);
        }

        if (hasReacted) {
          const filteredComments = [];
          for (let j = 0; j < msg.commenti.length; j++) {
            const comment = msg.commenti[j];
            if (!(comment.idautore == this.currentUserId && comment.reazione == reaction)) {
              filteredComments.push(comment);
            }
          }
          this.messages[i].commenti = filteredComments;
        } else {
          this.messages[i].commenti.push({
            idautore: this.currentUserId,
            reazione: reaction,
            commento_id: Date.now()
          });
        }
        break;
      }
    }

    this.$router.go();
    this.closeCommentMessageModal();

  } catch (error) {
    if (error.response) {
      const messaggio = error.response.data.errore;
      const codiceErrore = parseInt(error.response.data.codiceErrore);
      alert(messaggio + ' (codice ' + codiceErrore + ')');
    } else {
      alert("Errore: Errore di rete.");
    }
    this.$router.go();
   }
  },

  async loadMessages() {
    try {
      const response = await this.$axios.get(`/wasachat/${this.currentUser}/chats/${this.chatId}`);

      this.messages = [];
      for (let index = 0; index < response.data.length; index++) {
        const message = response.data[index];
        
        if (message.foto && !message.foto.startsWith("data:image")) {
          message.foto = 'data:image/jpeg;base64,' + message.foto;
        }
        
        if (!Array.isArray(message.commenti)) {
          message.commenti = [];
        }
        
        this.messages.push(message);
      }
      const self = this; 

    //auto_scroll
    this.$nextTick(function() {
      const container = self.$refs.messagesList;
      if (container) {
        container.scrollTop = container.scrollHeight;
      }

    });
  } catch (e) {
    if (e.response && e.response.data) {
      const messaggio = e.response.data.errore;
      const codiceErrore = parseInt(e.response.data.codiceErrore);
      alert(messaggio + ' (codice ' + codiceErrore + ')');
    }
    console.error(e);
    }
  },

  async deleteMessage() {
    try {
      const response = await this.$axios.delete(
        `/wasachat/${this.currentUser}/chats/${this.chatId}/messaggi/${this.selectedMessage.message_id}`
      );
        this.$router.go()

      this.closeMessageModal();
    } catch (error) {
      if (error.response) {
        const messaggio = error.response.data.errore;
        const codiceErrore = parseInt(error.response.data.codiceErrore);
        alert(messaggio + ` (codice ${codiceErrore})`);
      } else {
        alert("Error: Network error");
      }
    }
  },

  async sendMessage() {
    {
      const messageData = {
        testo: this.newMessage.trim(),
        foto: this.newPhoto || "",
      };
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
        
        const response =  await this.$axios.post(
          `/wasachat/${this.currentUser}/chats/${this.chatId}`,
          messageData
        );

        this.$router.go();
        
      } catch (error) {
      if (error.response) {
        const messaggio = error.response.data.errore;
        const codiceErrore = parseInt(error.response.data.codiceErrore);
        alert(messaggio + ` (codice ${codiceErrore})`);
      } else {
        alert("Error: Network error.");
      }
    }
   }
  },

  async sendReplyMessage() {
    const messageData = {
      testo: this.ans.trim(),
      foto: this.ansphoto || "",
    };

    try {
      const response = await this.$axios.post(
        `/wasachat/${this.currentUser}/risposta/chats/${this.chatId}/messaggi/${this.selectedMessage.message_id}`,
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
      if (error.response) {
        const messaggio = error.response.data.errore;
        const codiceErrore = parseInt(error.response.data.codiceErrore);
        alert(messaggio + ` (codice ${codiceErrore})`);
      } else {
        alert("Error: Network error.");
      }
    } finally {
      this.ans = "";
      this.ansphoto = "";
      this.closeanswereMessageModal();
    }
  },

  async showParticipants() {
    this.$router.push({ 
        name: 'GroupMembersView', 
        params: { nickname: this.currentUser, chat: this.$route.params.chat } 
      });
  },

  async goToForwardView() {
    this.$router.push({ 
        name: 'ForwardView', 
        params: { nickname: this.currentUser, chat: this.$route.params.chat, message: this.selectedMessage.message_id } 
      });
  },

  async quitGroup() {
    try {
      const response = await this.$axios.delete(`/wasachat/${this.currentUser}/chats/${this.chatId}`);
      this.$router.push({ 
        name: 'UserChats', 
        params: { nickname: this.currentUser } 
      });
            } catch (error) {
      if (error.response) {
        const messaggio = error.response.data.errore;
        const codiceErrore = parseInt(error.response.data.codiceErrore);
        alert(messaggio + ` (codice ${codiceErrore})`);
      } else {
        alert("Error: Network error.");
      }
    } finally {
      this.closeQuitModal();
    }
  },

  async addUserToGroup(nickname) {
    try {
      this.addingMember = true;
      
      const response = await this.$axios.put(
        `/wasachat/${this.currentUser}/chats/gruppi/${this.chatId}/aggiungi`,
        { utente_da_aggiungere: nickname.trim() }
      );
      
      alert(response.data.risposta);
      this.closeAddMemberModal();
      
    } catch (error) {
      if (error.response) {
        const messaggio = error.response.data.errore;
        const codiceErrore = parseInt(error.response.data.codiceErrore);
        alert(messaggio + ` (codice ${codiceErrore})`);
      } else {
        alert("Error: Network error.");
      }
    } finally {
      this.addingMember = false;
    }
   },
  }
};

</script>

<style scoped>
.messages_container_group {
  padding: 10px;
  max-width: 80%;
  margin: 0 auto;
  background-color: #ffffff;
  border-radius: 1px;
  padding-bottom: 2px;
  height: calc(100vh - 70px);
  overflow: hidden;
}

.messages_container_private {
  max-width: 90%;
  margin: 0 auto;
  background-color: #ffffff;
  border-radius: 1px;
  padding-bottom: 2px;
  height: calc(100vh - 1px);
  overflow: hidden;
}

.btn {
  background-color: rgb(125, 3, 240);
  color: rgb(255, 255, 255);
  border-radius: 90px;
  font-size: 15px;
  padding: 20px;
  margin: 10px;
  cursor: pointer;
}

.quit_btn {
  background-color: rgb(170, 86, 173);
  color: rgb(255, 255, 255);
  border-radius: 90px;
  font-size: 15px;
  padding: 20px;
  margin: 10px;
  cursor: pointer;
}

.btn_2{
  background-color: rgb(125, 3, 240);
  color: rgb(255, 255, 255);
  border-radius: 90px;
  padding: 12px;
  margin: 10px;
  font-size: 15px;
  cursor: pointer;
}

.messages_list {
  list-style-type: none;
  padding: 20px;
  margin: 10px;
  max-height: calc(100% - 120px);
  overflow-y: auto;
  scroll-behavior: smooth;
  width: 100%;
  box-sizing: border-box;
}

.message_item {
  display: flex;
  align-items: flex-start;
  padding: 16px;
  margin-bottom: 16px;
  background-color: rgb(220, 213, 228);
  border-radius: 10px;
  max-width: 60%;
  cursor: pointer
}

.message_item_right {
  display: flex;
  align-items: flex-start;
  padding: 16px;
  margin-bottom: 16px;
  margin-left: auto;
  background-color: rgb(209, 188, 230);
  border-radius: 10px;
  max-width: 60%;
  cursor: pointer
}

.message_item_group {
  border-radius: 10px;
  max-width: 60%;
  cursor: pointer;
  margin-right: auto;
  background-color: var(--random-group-color, rgb(220, 213, 228)); 
  display: flex;
  align-items: flex-start;
  padding: 16px;
  margin-bottom: 16px;
}

.message_sender {
  font-weight: bold;
  font-size: 20px;
  color: #171529;
  margin-bottom: 50px; 
  margin-right: 30px;
}

.message_photo-container {
  display: flex;
  flex-direction: column;
  margin-right: 30px;
  align-items: center;
}

.message_photo {
  margin-right: 30px;
  border-radius: 15px;
  margin-bottom: 15px;
  max-width: 75%;
  height: auto;
  max-height: 250px;
}

.message_text-container {
  display: flex;
  flex-direction: column;
}

.message_text {
  margin: 10;
  color: #495057;
  font-size: 16px;
  margin-top: 5px;
  margin-left: 15px;
  word-break: break-word;
  hyphens: auto;
}

.message_time {
  margin: 2px;
  color: #868e96;
  font-size: 12px;
}

.message_status {
  color: #0971d8;
  font-size: 14px;
  margin: 4px;
}

.message_reactions {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 8px;
  padding-top: 4px;
}

.reaction_badge_text {
  background-color: #f0f2f5;
  border-radius: 10px;
  padding: 2px 6px;
  font-size: 20px;
  cursor: pointer;
}

.reaction_badge_photo {
  background-color: #f0f2f5;
  border-radius: 10px;
  padding: 2px 6px;
  font-size: 20px;
  cursor: pointer;
}

.reaction_emoji_text {
  margin-right: 1px;
}

.reaction_emoji_photo {
  margin-right: 1px;
}

.message_input_group {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: rgb(209, 188, 230);
  border-top: 1px solid #e0e0e0;
  width: 95%;
  border-radius: 20px;
  margin: 0 auto;
  margin-bottom: 10px;
}

.message_input_private {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: rgb(209, 188, 230);
  border-top: 1px solid #e0e0e0;
  width: 95%;
  height: 120px;
  border-radius: 20px;
  margin: 0 auto;
  margin-bottom: 10px;
}

.message_input {
  flex-grow: 1;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  margin-right: 10px;
  font-size: 16px;
}

.group_header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  background-color: rgb(209, 188, 230);
  border-radius: 20px;
}

.group_actions {
  display: flex;
  gap: 10px;
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

.modal_content {
  background-color: white;
  border-radius: 10px;
  padding: 20px;
  width: 90%;
  max-width: 400px;
}

.modal_content-large {
  background-color: white;
  border-radius: 10px;
  padding: 30px;
  max-width:1250px;
  max-height: 1000px;
}

.modal_input {
  width: 100%;
  padding: 10px;
  border-radius: 20px;
  margin: 10px 0;
  font-size: 16px;
}

.modal_btn {
  padding: 8px 15px;
  background-color: rgb(125, 3, 240);
  color:rgb(255, 255, 255);
  border-radius: 20px;
  cursor: pointer;
  gap: 10px;
  margin-left: 10px;
  font-size: 14px;
  margin: 10px 0;
  margin-right: 10px;
}

.modal_btn_gray {
  border-radius: 20px;
  cursor: pointer;
  gap: 10px;
  margin-left: 10px;
  font-size: 14px;
  margin: 10px 0;
  margin-right: 10px;
  padding: 8px 15px;
  background-color: #6c757d;
  color: white;
}

.modal_btn_red{
  border-radius: 20px;
  cursor: pointer;
  gap: 10px;
  margin-left: 10px;
  font-size: 14px;
  margin: 10px 0;
  margin-right: 10px;
  padding: 8px 15px;
  background-color: rgb(161, 63, 84);
  color: white;
}

.message_reply-container {
  border-left: 3px solid #000000;
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
  color: #000000;
  margin-right: 10px;

}

.reaction_button {
  font-size: 24px;
  background: none;
  cursor: pointer;
}

.forward_label {
  font-size: 10px;
  color: #666;
  margin-bottom: 5px;
  padding: 2px 6px;
  border: 1px solid #666;
  border-radius: 90px;
}

.message_input_ans {
  flex-grow: 1;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  margin-right: 10px;
  font-size: 16px;
  width: 100%;
  margin-bottom: 10px;
}

.no_messages {
  max-width: 700px;
  background-color: rgb(209, 188, 230);
  margin: 180px auto;
  font-size: 60px;
  padding: 80px;
  border-radius: 20px
}

.no_conversations {
  max-width: 700px;
  background-color: rgb(209, 188, 230);
  margin: 180px auto;
  font-size: 60px;
  padding: 80px;
  border-radius: 20px
}

.goBack_btn {
  background-color: rgb(161, 63, 84);
  color: rgb(221, 219, 219);
  padding: 20px 40px;
  margin: 40px;
  border-radius: 90px;
  font-size: 15px;
  position: fixed;
  top: 0px;    
  right: 10px;      
}
</style>