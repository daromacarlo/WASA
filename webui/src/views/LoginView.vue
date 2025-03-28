
<template>
    <div class="login-container">
      <h2>Accedi</h2>
      <form @submit.prevent="login">
        <input type="text" v-model="nickname" placeholder="Nickname" required />
        <button type="submit" :disabled="loading">Accedi</button>
      </form>
      <p v-if="errormsg" class="error">{{ errormsg }}</p>
      <p v-if="loading">Caricamento...</p>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        nickname: '', 
        errormsg: null,
        loading: false,
        some_data: null,
      };
    },
    methods: {
  async login() {
    if (!this.nickname.trim()) {
      this.errormsg = "Inserisci un nickname valido!";
      return;
    }

    this.loading = true;
    this.errormsg = null;

    try {
      // Effettua la richiesta di login
      const response = await this.$axios.post("/wasachat", { nickname: this.nickname });
      
      if (response.status === 200) {
        // Se il login è andato a buon fine, naviga alla pagina delle conversazioni
        this.$router.push(`/wasachat/${this.nickname}/chats`);  // Passa il nickname come parametro nella URL
      } else {
        this.errormsg = "Errore durante il login. Riprova.";
      }
    } catch (e) {
      // Gestisci gli errori (es. utente non trovato)
      if (e.response && e.response.status === 404) {
        this.errormsg = "Utente non trovato. Verifica il nickname.";
      } else {
        this.errormsg = "Errore nel tentativo di login. Riprova.";
      }
      console.error(e);
    } finally {
      this.loading = false;
    }
  },
    },
  };
  </script>
  
  <style scoped>
  .login-container {
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
  
  .error {
    color: red;
    font-size: 14px;
  }
  </style>
  