<template>
  <div class="row justify-center">
    <div class="column q-pa-lg">
      <div class="row">
        <q-card square class="shadow-24" style="width:400px;">
          <q-card-section>
            <h4>Keeper</h4>
          </q-card-section>
          <q-form class="q-px-sm q-pt-xl" @submit="submit">
            <q-card-section>
              <q-input  
                square 
                clearable 
                v-model="password"                                                        
                type="password"  
                lazy-rules
                label="Password">
              </q-input>
            </q-card-section>
            <q-card-actions class="q-px-lg">
              <q-btn 
                size="lg" 
                color="secondary"
                class="full-width text-white"
                type="submit"                      
                label="Login" />
            </q-card-actions>
          </q-form>
        </q-card>
      </div>
    </div>
  </div>

</template>

<script setup lang="ts">
import { useTokenStore } from 'src/stores/token-store';
import { ref } from 'vue';
import { useRouter } from 'vue-router'

const tokenStore = useTokenStore();
const password = ref('');
const router = useRouter();

const submit = async (e: Event | SubmitEvent) => {
  e.preventDefault();
  await tokenStore.login(password.value);
  router.push('/');
  return false;
}

</script>
