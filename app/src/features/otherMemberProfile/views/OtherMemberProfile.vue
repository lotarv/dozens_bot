<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { useMembersStore } from '@/stores/membersStore'
import { Member } from '@/types/Member'
import ArrowLeft from '@/components/icons/ArrowLeft.vue'
import Reports from '../ui/Reports.vue'
import Declarations from '../ui/Declarations.vue'
import { useReportsStore } from '@/stores/reportsStore'
import { useUserStore } from '@/stores/user'
import { onBeforeMount, ref } from 'vue'
import { DeclarationDocument } from '../entities/DeclarationDocument'
import { DozensTransport } from '@/repository/http'
const route = useRoute()
const membersStore = useMembersStore()
const username = route.params.username.toString()


const loaded = ref<boolean>(false)
const reportsStore = useReportsStore()
const userStore = useUserStore()

const declarations = ref<DeclarationDocument[]>([])

const member: Member = membersStore.getMemberByUsername(username)



onBeforeMount(async () => {
    await reportsStore.fetchUserReports(username)
    loaded.value = true
    if(declarations.value.length == 0) declarations.value = await DozensTransport.GetDeclarations(username)
})
</script>
<template>
    <section class="absolute inset-0 z-10 overflow-y-auto">
        <div class="header" :style="{backgroundImage: `url(${member.avatar_url})`}">
            <div class="header-link-and-income">
                <RouterLink :to="'/members'">
                    <div class="p-2 bg-white w-12 h-12 rounded-full text-[22px] flex items-center justify-center">
                        <ArrowLeft />
                    </div>
                </RouterLink>

                <div class="member-income">
                    <span>{{ `${member.annual_income }М ₽ / год`}}</span>
                </div>

            </div>
            <div class="blur-name-box">
                <span class="font-medium text-4xl tracking-[-1px]">{{ member.fio }}</span>
                <span class="font-medium text-base tracking-[-0.4px]">{{ member.niche }}</span>
            </div>
        </div>
        <div class="documents -mt-4 z-20" v-if="loaded">
            <div class="declarations-container" v-if="declarations && declarations.length > 0">
                <Declarations :declarations="declarations" :member="member"/>
            </div>
            <div class="reports-container">
                <Reports v-if="reportsStore.reports[username]" :reports="reportsStore.reports[username].reports ?? []" :username="username" />
            </div>
        </div>
    </section>
</template>

<style scoped>
.reset-layout {
  margin: 0 !important;
  padding: 0 !important;
}

section {
    @apply flex flex-col; 
}

.header{
    @apply flex flex-col gap-[12rem] bg-cover bg-center 
}

.header-link-and-income{
    font-family: 'SF Pro Text';
    @apply flex justify-between p-3
}

.blur-name-box {
    @apply flex flex-col gap-2 text-white pt-4 px-5 pb-7;

}

.blur-name-box {
  position: relative;
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  color: white;
  overflow: hidden;
}

/* Градиент под блюром */
.blur-name-box::before {
  content: "";
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.8), rgba(0, 0, 0, 0));
  z-index: 0;
}

/* Контент поверх градиента */
.blur-name-box > * {
  position: relative;
  z-index: 1;
}


.member-income{
    @apply flex flex-row gap-2 text-[16px] font-semibold tracking-[-0.4px] leading-6 items-center rounded-full bg-[hsla(0,0%,100%,0.7)] px-3 py-2;
    backdrop-filter:blur(20.1px);
    -webkit-backdrop-filter: blur(20.1px); /* для Safari */
}

.member-name {
    font-family: "SF Pro Display";
    @apply font-[500] text-4xl tracking-[-1px] leading-9 pl-[9px]
}
.dot {
    @apply w-[4px] h-[4px] bg-black inline-block rounded-[50%] self-center
}
</style>