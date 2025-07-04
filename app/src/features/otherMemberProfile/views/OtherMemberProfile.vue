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
    declarations.value = await DozensTransport.GetDeclarations(username)
    loaded.value = true
})
</script>
<template>
    <section>
        <div class="header">
            <div class="header-link-and-income">
                <RouterLink :to="'/members'">
                    <div class="p-2 bg-white w-12 h-12 rounded-full text-[22px] flex items-center justify-center">
                        <ArrowLeft />
                    </div>
                </RouterLink>

                <div class="member-income">
                    <span>{{ `${member.annual_income }М ₽ / год`}}</span>
                    <div class="dot"></div>
                    <span class="niche">{{ member.niche }}</span>
                </div>

            </div>
            <div class="member-name">
                {{ member.fio }}
            </div>
        </div>
        <div class="documents" v-if="loaded">
            <div class="declarations-container">
                <Declarations :declarations="declarations" :member="member"/>
            </div>
            <div class="reports-container">
                <Reports v-if="reportsStore.reports[username]" :reports="reportsStore.reports[username].reports ?? []" :username="username" />
            </div>
        </div>
    </section>
</template>

<style scoped>
section {
    @apply flex flex-col pl-1 pb-1 pr-1 pt-[7px]; 
}

.header{
    @apply flex flex-col gap-4 p-3
}

.header-link-and-income{
    font-family: 'SF Pro Text';
    @apply flex justify-between 
}

.member-income{
    @apply flex flex-row gap-2 text-[16px] font-[600] tracking-[-0.4px] leading-6 items-center
}

.member-name {
    font-family: "SF Pro Display";
    @apply font-[500] text-4xl tracking-[-1px] leading-9 pl-[9px]
}
.dot {
    @apply w-[4px] h-[4px] bg-black inline-block rounded-[50%] self-center
}
</style>