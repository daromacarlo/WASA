<template>

    <button @click="goBack" class="cancel-button">
      Indietro
    </button>

    <div class="modifyGroup">
      <div class="container">
        <h1 class="title">Modifica gruppo</h1>
        <div class="button-container">
          <button @click="openModifyNameModal" class="group-action-button" title="Modifica nome">
            <span class="button-text">Modifica nome</span>
          </button>
          <button @click="openModifyPhotoModal" class="group-action-button" title="Modifica foto">
            <span class="button-text">Modifica foto</span>
          </button>
        </div>
      </div>
  
      <!-- Modale modifica nome -->
      <div v-if="showModifyNameModal" class="modal">
        <div class="modal-content">
          <h3>Modifica nome gruppo</h3>
          <input 
            v-model="newGroupName" 
            type="text" 
            placeholder="Inserisci il nuovo nome del gruppo" 
            class="modal-input"
            @keyup.enter="modifyGroupName"
            :disabled="isProcessing"
          />
          <div class="modal-buttons">
            <button 
              @click="modifyGroupName" 
              class="modal-button confirm"
              :disabled="!newGroupName || isProcessing"
            >
              <span v-if="isProcessing">Salvando...</span>
              <span v-else>Salva</span>
            </button>
            <button @click="closeModifyNameModal" class="modal-button cancel">Annulla</button>
          </div>
        </div>
      </div>
  
      <!-- Modale modifica foto -->
      <div v-if="showModifyPhotoModal" class="modal">
        <div class="modal-content">
          <h3>Modifica foto gruppo</h3>
          <input 
            type="file" 
            accept="image/jpeg" 
            @change="handleFileUpload" 
            class="modal-input"
            :disabled="isProcessing"
          />
          <div v-if="errorMsg" class="error-message">{{ errorMsg }}</div>
          <div class="modal-buttons">
            <button 
              @click="modifyGroupPhoto" 
              class="modal-button confirm"
              :disabled="!newGroupPhoto || isProcessing"
            >
              <span v-if="isProcessing">Caricando...</span>
              <span v-else>Salva foto</span>
            </button>
            <button @click="closeModifyPhotoModal" class="modal-button cancel">Annulla</button>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        showModifyNameModal: false,
        showModifyPhotoModal: false,
        newGroupName: "",
        newGroupPhoto: "",
        isProcessing: false,
        errorMsg: null
      };
    },
    methods: {
      // Metodi per la modifica del nome
      openModifyNameModal() {
        this.showModifyNameModal = true;
        this.newGroupName = "";
        this.isProcessing = false;
      },
  
      closeModifyNameModal() {
        this.showModifyNameModal = false;
        this.newGroupName = "";
        this.isProcessing = false;
      },
  
      async modifyGroupName() {
        if (!this.newGroupName || !this.newGroupName.trim()) {
          alert("Inserisci un nome valido per il gruppo");
          return;
        }
  
        try {
          this.isProcessing = true;
          const currentUser = this.$route.params.nickname;
          const chatId = this.$route.params.chat;
          
          await this.$axios.put(
            `/wasachat/${currentUser}/gruppi/${chatId}/nome`,
            { nome: this.newGroupName.trim() }
          );
          
          alert("Nome del gruppo modificato con successo!");
          this.closeModifyNameModal();
          // Ricarica i dati del gruppo se necessario
          // await this.loadGroupInfo();
          
        } catch (error) {
          console.error("Errore nella modifica del nome:", error);
          this.handleGroupError(error);
        } finally {
          this.isProcessing = false;
        }
      },
  
      // Metodi per la modifica della foto
      handleFileUpload(event) {
        const file = event.target.files[0];
        this.errorMsg = null;
        
        if (!file) return;
        
        // Verifica che il file sia un'immagine JPEG
        if (!file.type.match(/image\/jpeg/)) {
          this.errorMsg = "Seleziona un'immagine in formato JPEG (.jpg o .jpeg)";
          this.newGroupPhoto = '';
          return;
        }
  
        this.convertToBase64(file);
      },
  
      convertToBase64(file) {
        const reader = new FileReader();
        reader.onload = () => {
          this.newGroupPhoto = reader.result;
        };
        reader.onerror = (error) => {
          console.error("Errore durante la conversione:", error);
          this.errorMsg = "Errore durante il caricamento dell'immagine";
        };
        reader.readAsDataURL(file);
      },
  
      openModifyPhotoModal() {
        this.showModifyPhotoModal = true;
        this.newGroupPhoto = "";
        this.errorMsg = null;
        this.isProcessing = false;
      },
  
      closeModifyPhotoModal() {
        this.showModifyPhotoModal = false;
        this.newGroupPhoto = "";
        this.errorMsg = null;
        this.isProcessing = false;
      },
  
      async modifyGroupPhoto() {
        if (!this.newGroupPhoto) {
          this.errorMsg = "Seleziona un'immagine valida";
          return;
        }
  
        try {
          this.isProcessing = true;
          const currentUser = this.$route.params.nickname;
          const chatId = this.$route.params.chat;
          
          await this.$axios.put(
            `/wasachat/${currentUser}/gruppi/${chatId}/foto`,
            { foto: this.newGroupPhoto }
          );
          
          alert("Foto del gruppo modificata con successo!");
          this.closeModifyPhotoModal();
          // Ricarica i dati del gruppo se necessario
          // await this.loadGroupInfo();
          
        } catch (error) {
          console.error("Errore nella modifica della foto:", error);
          this.handleGroupError(error);
        } finally {
          this.isProcessing = false;
        }
      },
  
      // Gestione errori comune
      handleGroupError(error) {
        let errorMessage = "Si è verificato un errore durante la modifica";
        if (error.response) {
          switch(error.response.status) {
            case 400: errorMessage = "Richiesta malformata"; break;
            case 401: errorMessage = "Non autorizzato"; break;
            case 403: errorMessage = "Non hai i permessi per questa modifica"; break;
            case 404: errorMessage = "Gruppo non trovato"; break;
            case 409: errorMessage = "Modifica non consentita"; break;
            case 413: errorMessage = "L'immagine è troppo grande"; break;
            case 500: errorMessage = "Errore del server. Riprova più tardi"; break;
          }
        } else if (error.request) {
          errorMessage = "Impossibile connettersi al server. Verifica la tua connessione";
        }
        
        alert(errorMessage);
      },
      goBack() {
      const { nickname, chat } = this.$route.params;
      this.$router.push(`/wasachat/${nickname}/chats/${chat}`);
    },
    }
  };
  </script>
  
  <style scoped>
  .modifyGroup {
    background-color: #f7f7f7;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    font-family: 'Arial', sans-serif;
  }
  
  .container {
    text-align: center;
    background-color: white;
    padding: 30px;
    border-radius: 10px;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
    width: 80%;
    max-width: 600px;
  }
  
  .title {
    font-size: 40px;
    color: #333;
    margin-bottom: 30px;
  }
  
  .button-container {
    display: flex;
    flex-direction: column;
    gap: 15px;
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
  
  .group-action-button {
    padding: 12px 20px;
    border: none;
    border-radius: 5px;
    font-size: 16px;
    cursor: pointer;
    transition: background-color 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    background-color: #7dac10;
    color: white;
    width: 100%;
  }
  
  .group-action-button:hover {
    background-color: #6a950d;
  }
  
  /* Stili per le modali */
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
    padding: 25px;
    border-radius: 8px;
    width: 90%;
    max-width: 400px;
  }
  
  .modal-input {
    width: 100%;
    padding: 10px;
    margin: 15px 0;
    border: 1px solid #ddd;
    border-radius: 4px;
  }
  
  .modal-buttons {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    margin-top: 20px;
  }
  
  .modal-button {
    padding: 8px 15px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  .confirm {
    background-color: #7dac10;
    color: white;
  }
  
  .cancel {
    background-color: #f0f0f0;
    color: #333;
  }
  
  .error-message {
    color: #d32f2f;
    margin: 10px 0;
    font-size: 14px;
  }

  </style>