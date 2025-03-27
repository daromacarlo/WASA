<template>
  <div class="register-container">
    <h2>Registrati</h2>
    <form @submit.prevent="register">
      <input type="text" v-model="nickname" placeholder="Nickname" required />
      <!-- Input per selezionare l'immagine (solo JPEG) -->
      <input type="file" @change="handleFileUpload" accept="image/jpeg" required />
      <button type="submit" :disabled="loading">
        {{ loading ? "Caricamento..." : "Registrati" }}
      </button>
    </form>
    <p v-if="errorMsg" class="error-message">{{ errorMsg }}</p>
  </div>
</template>

<script>
export default {
  data() {
    return {
      nickname: '',
      fotoProfilo: '', // Qui memorizzeremo la stringa Base64
      errorMsg: null,
      loading: false,
    };
  },
  methods: {
    // Gestisce la selezione del file
    handleFileUpload(event) {
      const file = event.target.files[0]; // Ottieni il file selezionato
      if (file) {
        // Verifica che il file sia un'immagine JPEG
        if (!file.type.match(/image\/jpeg/)) {
          this.errorMsg = "Seleziona un'immagine in formato JPEG (.jpg o .jpeg).";
          this.fotoProfilo = ''; // Resetta la foto profilo
          return;
        }
        this.convertToBase64(file); // Converti il file in Base64
      } else {
        this.errorMsg = "Seleziona un'immagine valida!";
      }
    },

    // Converte il file in una stringa Base64
    convertToBase64(file) {
      const reader = new FileReader();
      reader.onload = () => {
        this.fotoProfilo = reader.result; // Salva la stringa Base64
        this.errorMsg = null; // Resetta il messaggio di errore
      };
      reader.onerror = (error) => {
        console.error("Errore durante la conversione del file:", error);
        this.errorMsg = "Errore durante il caricamento dell'immagine.";
      };
      reader.readAsDataURL(file); // Avvia la conversione
    },

    // Invia i dati al backend
    async register() {
      if (!this.nickname.trim()) {
        this.errorMsg = "Inserisci un nickname valido!";
        return;
      }

      if (!this.fotoProfilo) {
        this.errorMsg = "Seleziona un'immagine!";
        return;
      }

      this.loading = true;
      this.errorMsg = null;

      try {
        const response = await this.$axios.put("/wasachat", {
          nickname: this.nickname,
          foto: this.fotoProfilo, // Invia la stringa Base64
        });

        // Se la risposta ha successo (status 2xx), considera la registrazione riuscita
        if (response.status >= 200 && response.status < 300) {
          alert("Registrazione effettuata con successo!");
          this.$router.push(`/`);
        } else {
          this.errorMsg = "Errore durante la registrazione. Riprova.";
        }
      } catch (e) {
        console.error("Errore durante la registrazione:", e);
        // Se il backend restituisce un messaggio di errore, usalo
        this.errorMsg = e.response?.data?.message || "Errore durante la registrazione. Riprova.";
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<style scoped>
.register-container {
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