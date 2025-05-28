<template>
  <button @click="goBack" class="goBack_btn">Go Back</button>
    <div class="c">
      <h1>Modify User Profile</h1>
      <div class="bc">
        <button @click="openModifyNameModal" class="btn" title="Insert your new nickname">Modify Nickname</button>
        <button @click="openModifyPhotoModal" class="btn" title="Insert your new photo">Modify Photo</button>
    </div>

    <div v-if="showModifyNameModal" class="modal">
      <div class="modal-content">
        <h3>Modify your Nickname</h3>
        <form @submit.prevent="modifyGroupName">
          <input v-model="newNick" type="text" placeholder="Insert your new name:" class="modal-input" required/>
          <div class="modal-btn">
            <button class="btn">Save</button>
            <button @click="closeModifyNameModal" class="btn_gray">Go Back</button>
          </div>
        </form>
      </div>
    </div>

    <div v-if="showModifyPhotoModal" class="modal">
      <div class="modal-content">
        <h3>Modify your photo</h3>
        <form @submit.prevent="modifyGroupPhoto">
          <input type="file" accept="image/jpeg" class="modal-input" @change="handleFileUpload" required/>
          <div class="modal-btn">
            <button class="btn">Save</button>
            <button @click="closeModifyPhotoModal" class="btn_gray">Go Back</button>
          </div>
        </form>
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
    newNick: "",
    newPhoto: "",
    error: null,
    currentUser: this.$route.params.nickname
   };
  },

  methods: {

  handleFileUpload(event) {
    const file = event.target.files[0];
    this.error = null;
    if (!file) return;
    if (!file.type.match(/image\/jpeg/)) {
      this.error = "Select a jpg or jpeg image";
      this.newPhoto = '';
      return;
    }
    this.convertToBase64(file);
  },
  
  convertToBase64(file) {
    const reader = new FileReader();
    reader.onload = () => {
      this.newPhoto = reader.result;
    };
    reader.onerror = (error) => {
      console.error("Error:", error);
      this.error = "Error: error during loading image";
    };
    reader.readAsDataURL(file);
  },

  openModifyNameModal() {
    this.showModifyNameModal = true;
    this.newNick = "";
  },

  closeModifyNameModal() {
    this.showModifyNameModal = false;
    this.newNick = "";
  },

  openModifyPhotoModal() {
    this.showModifyPhotoModal = true;
    this.newPhoto = "";
    this.error = null;
  },

  closeModifyPhotoModal() {
    this.showModifyPhotoModal = false;
    this.newPhoto = "";
    this.error = null;
  },

  async modifyGroupName() {
    try {
      const newName = this.newNick.trim() 
      const response = await this.$axios.put(
        `/wasachat/${this.currentUser}/nome`,
        { nome: this.newNick.trim() }
      );

      const codice = parseInt(response.data.codice);
      const message = response.data.risposta || "";

      if (codice >= 200 && codice < 300) {
        alert(message)
        this.closeModifyNameModal();
        this.$router.push({
          name: 'UserChats',
          params: { nickname:newName }
        });
      } else {
        alert(message);
      }

    } catch (e) {
      if (e.response && e.response.data) {
        const message = e.response.data.errore;
        const codiceErrore = e.response.data.codiceErrore
        parseInt(e.response.data.codiceErrore)
        alert(`${message} (codice ${codiceErrore})`);
      } else {
        alert("Error: Network error.");
      }
    } finally {
      console.error(e);
    }
  },

  async modifyGroupPhoto() {
    try {
      const response = await this.$axios.put(
        `/wasachat/${this.currentUser}/foto`,
        { foto: this.newPhoto }
      );
      const message = response.data.risposta;
      const codice = parseInt(response.data.codice);

      this.$router.push({ 
      name: 'UserChats', 
      params: { nickname: this.currentUser } 
    });
    if (codice == 200) {
      alert(message);
      this.closeModifyPhotoModal();
    }
  } catch (e) {
    if (e.response) {
      const message = e.response.data.errore;
      const codiceErrore = parseInt(e.response.data.codiceErrore);
      alert(message + ` (codice ${codiceErrore})`);
    } else {
      alert("Error: Network error.");
    }
  }
  finally {
  console.error(e);
  }
},
  
  goBack() {
    this.$router.push(`/wasachat/${this.currentUser}/chats`);
    },
  },
};
  </script>
  
  <style scoped>

  .goBack_btn {
    background-color: rgb(161, 63, 84);
    color: rgb(221, 219, 219);
    padding: 20px 40px;
    margin: 40px;
    border-radius: 90px;
    font-size: 15px;
    box-shadow: 0 6px 15px rgba(0, 0, 0, 0.1);
    position: fixed;
    top: 0px;
    right: 40px;
    border: none;
    cursor: pointer;
  }
  
  .c {
    text-align: center;
    background-color: rgb(209, 188, 230);
    padding: 40px;
    border-radius: 12px;
    width: 60%;
    margin: auto;
    border-radius: 5px;
    margin-top: 400px;
  }
  
  .bc {
    display: flex;
    flex-direction: column;
    gap: 15px;
    padding: 40px;
  }

  .btn {
    background-color: rgb(125, 3, 240);
    color: rgb(255, 255, 255);
    padding: 20px 40px;
    border-radius: 90px;
    font-size: 15px;
    margin-right: 30px;
    margin-left: 30px;
    border: none;
    cursor: pointer;
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
    padding: 25px;
    border-radius: 8px;
    width: 90%;
    max-width: 700px;
  }
  
  .modal-input {
    width: 100%;
    padding: 10px;
    margin: 15px 0;
    border: 1px solid #ddd;
    border-radius: 4px;
  }
  
  .btn_gray {
    background-color: rgb(172, 159, 184);
    color: rgb(255, 255, 255);
    padding: 20px 40px;
    border-radius: 90px;
    font-size: 15px;
    margin-right: 30px;
    margin-left: 30px;
    border: none;
    cursor: pointer;
  }

  </style>