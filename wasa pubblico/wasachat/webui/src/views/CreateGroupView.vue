<template>
  <div class="creategroup-container">
    <h2>Crea il tuo gruppo</h2>
    <form @submit.prevent="createGroup">
      <input type="text" v-model="nome" placeholder="Nome Gruppo" required />
      <!-- Input per selezionare l'immagine (solo JPEG) -->
      <input type="file" @change="handleFileUpload" accept="image/jpeg" required />
      <button type="submit" :disabled="loading">
        {{ loading ? "Caricamento..." : "Crea gruppo" }}
      </button>
    </form>
    <p v-if="errorMsg" class="error-message">{{ errorMsg }}</p>
  </div>
</template>

<script>
export default {
  data() {
    return {
      nome: '', // Nome del gruppo
      foto: '', // Stringa Base64 dell'immagine
      errorMsg: null, // Messaggio di errore
      loading: false, // Stato di caricamento
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
          this.foto = ''; // Resetta la foto
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
        this.foto = reader.result; // Salva la stringa Base64
        this.errorMsg = null; // Resetta il messaggio di errore
      };
      reader.onerror = (error) => {
        console.error("Errore durante la conversione del file:", error);
        this.errorMsg = "Errore durante il caricamento dell'immagine.";
      };
      reader.readAsDataURL(file); // Avvia la conversione
    },

    // Invia i dati al backend per creare il gruppo
    async createGroup() {
      if (!this.nome.trim()) {
        this.errorMsg = "Inserisci un nome valido per il gruppo!";
        return;
      }

      if (!this.foto) {
        this.errorMsg = "Seleziona un'immagine!";
        return;
      }

      this.loading = true;
      this.errorMsg = null;

      try {
        const nickname = this.$route.params.nickname; // Ottieni il nickname dall'URL
        const response = await this.$axios.post(`/wasachat/${nickname}/gruppi`, {
          nome: this.nome,
          foto: this.foto, // Invia la stringa Base64
        });

        // Se la risposta ha successo (status 2xx), reindirizza alla pagina delle chat
        if (response.status >= 200 && response.status < 300) {
          alert("Gruppo creato con successo!");
          this.$router.push(`/wasachat/${nickname}/chats`);
        } else {
          this.errorMsg = "Errore durante la creazione del gruppo. Riprova.";
        }
      } catch (e) {
        console.error("Errore durante la creazione del gruppo:", e);
        this.errorMsg = e.response?.data?.message || "Errore durante la creazione del gruppo. Riprova.";
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<style scoped>
.creategroup-container {
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