<template>
    <div class="searchuser-container">
      <h2>Cerca una persona</h2>
      <form @submit.prevent="searchUser">
        <input type="text" v-model="utente" placeholder="Nome utente" required />
        <button type="submit" :disabled="loading">
          {{ loading ? "Caricamento..." : "Cerca" }}
        </button>
      </form>
      <p v-if="errorMsg" class="error-message">{{ errorMsg }}</p>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        utente: '', // Nome utente da cercare
        errorMsg: null, // Messaggio di errore
        loading: false, // Stato di caricamento
      };
    },
    methods: {
      // Cerca un utente e crea una chat privata
      async searchUser() {
        if (!this.utente.trim()) {
          this.errorMsg = "Inserisci un nome utente valido!";
          return;
        }
  
        this.loading = true;
        this.errorMsg = null;
  
        try {
          const nickname = this.$route.params.nickname; // Ottieni il nickname dall'URL
          const response = await this.$axios.post(`/wasachat/${nickname}/conversazioniprivate`, {
            utente: this.utente, // Invia il nome utente da cercare
          });
  
          // Se la risposta ha successo (status 2xx), reindirizza alla pagina delle chat
          if (response.status >= 200 && response.status < 300) {
            alert("Chat privata creata con successo!");
            this.$router.push(`/wasachat/${nickname}/chats`);
          } else {
            this.errorMsg = "Errore durante la creazione della chat privata. Riprova.";
          }
        } catch (e) {
          console.error("Errore durante la creazione della chat privata:", e);
          this.errorMsg = e.response?.data?.message || "Errore durante la creazione della chat privata. Riprova.";
        } finally {
          this.loading = false;
        }
      },
    },
  };
  </script>
  
  <style scoped>
  .searchuser-container {
    padding: 20px;
    max-width: 400px;
    margin: 0 auto;
    text-align: center;
  }
  
  input {
    display: block;
    width: 100%;
    padding: 10px;
    margin: 10px 0;
    border: 1px solid #ccc;
    border-radius: 5px;
  }
  
  button {
    padding: 10px 20px;
    background-color: #7dac10;
    color: white;
    border: none;
    border-radius: 5px;
    cursor: pointer;
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