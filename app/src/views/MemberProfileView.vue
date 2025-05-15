<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { useMembersStore } from '@/stores/membersStore'
import { Member } from '@/types/Member'
import ArrowLeft from '@/components/icons/ArrowLeft.vue'
import { declarations } from '@/mocks/declarations'
import DeclarationsComponent from '@/components/memberProfileView/DeclarationsComponent.vue'
import ReportsComponent from '@/components/memberProfileView/ReportsComponent.vue'
import { useReportsStore } from '@/stores/reportsStore'
import { useUserStore } from '@/stores/user'
import { onBeforeMount } from 'vue'
const route = useRoute()
const membersStore = useMembersStore()
const username = route.params.username.toString()

const reportsStore = useReportsStore()
const userStore = useUserStore()


const member: Member = membersStore.getMemberByUsername(username)

onBeforeMount(async () => {
    reportsStore.fetchUserReports(username)
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
        <div class="declarations-container">
            <DeclarationsComponent :declarations="declarations"/>
        </div>
        <div class="reports-container">
            <ReportsComponent v-if="reportsStore.reports[username]" :reports="reportsStore.reports[username].reports ?? []" :username="username" />
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
    @apply font-[500] text-4xl tracking-[-1px] leading-9
}
.dot {
    @apply w-[4px] h-[4px] bg-black inline-block rounded-[50%] self-center
}
</style>