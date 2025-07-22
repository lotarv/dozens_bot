<script lang="ts" setup>
import axios from 'axios';
import { ref, onBeforeMount, Ref } from 'vue';
import { getTelegramInitData } from '../services/auth';
import { current_user } from '@/mocks/current_user';
import { meetings } from '@/mocks/meetings';
import BellIcon from '@/components/icons/BellIcon.vue';
import MeetingsSlider from '@/components/mainView/MeetingsSlider.vue';
import ArrowIcon from '@/components/icons/ArrowIcon.vue';
import ArrowUp from '@/components/icons/ArrowUp.vue';
import MembersAvatars from '@/components/mainView/MembersAvatars.vue';
import MeetingCard from '@/components/mainView/MeetingCard.vue';
// import { members } from '@/mocks/members';
import SadFaceIcon from '@/components/icons/SadFaceIcon.vue';
import { declaration } from '@/mocks/declaration';
import CurrentDeclaration from '@/components/mainView/CurrentDeclaration.vue';
import { useMembersStore } from '@/stores/membersStore';
import { useDecryptionStore } from '@/stores/decryption';
import { useRouter } from 'vue-router';
import { useBankStore } from '@/features/piggyBank/model/bankStore';
interface Member {
    fio: string;
    avatar_url: string;
    niche: string;
    annual_income: number;
    username: string;
}

const membersStore = useMembersStore()
const cryptStore = useDecryptionStore()
const bankStore = useBankStore()
const isLoading = ref(false);
const error = ref<string | null>(null);
const router = useRouter()

function shortenLastName(fullname: string): string {
    let separated = fullname.split(" ")
    if (separated.length == 1) {
        return fullname
    } else {
        return separated[0] + " " + separated[1][0] + "."
    }
}
//TODO: add user store
async function createOrUpdateUser() {
    try {
        await axios.post(
            `${import.meta.env.VITE_API_URL}/users`,
            {},
            {
                headers: {
                    'X-Telegram-Init-Data': getTelegramInitData(),
                },
            }
        );
    } catch (err) {
        error.value = 'Failed to authenticate. Please try again.';
        console.error('Authentication failed:', err);
    }
}
const isReady = ref(false)
onBeforeMount(async () => {
    window.Telegram.WebApp.disableVerticalSwipes();
    await createOrUpdateUser();
    await cryptStore.fetchKey();
    if (cryptStore.key == "") {
        router.push({name: "login"})
    }
    await membersStore.fetchMembers();
    await bankStore.fetchPiggyBank()
    isReady.value = true

});
</script>

<template>
    <section class="flex flex-col p-1 gap-1" v-if="isReady">
        <div class="header">
            <div class="cur-user-name">{{ shortenLastName(current_user.fio) }}</div>
            <div class="cur-user-benefits">
                <img :src="current_user.avatar_url" alt="your_pic" class="cur-user-pic">
                <div class="notifications">
                    <BellIcon />
                </div>
            </div>
        </div>
        <div class="slider-container">
            <MeetingsSlider :meetings="meetings"/>
        </div>
        <div class="info-block-1">
            <div class="flex-1 inactive">
                <CurrentDeclaration :declaration="declaration"/>
            </div>
            <div class="flex-1 border-3 border-black">
                <div class="members-and-bank">
                <RouterLink :to="'/members'">
                    <div class="block members">
                        <div class="block-header">
                            <div class="block-title">Участники</div>
                            <div class="text-[24px]">
                                <ArrowIcon />
                            </div>
                        </div>
                        <MembersAvatars :current_user="current_user" :members="membersStore.members" />
                    </div>
                </RouterLink>
                <div class="block" @click="router.push('/piggy-bank')">
                    <div class="block-header">
                        <div class="block-title">Копилка</div>
                        <div class="text-[24px]">
                            <ArrowIcon />
                        </div>
                    </div>
                    <div class="block-info">
                        <div class="status">
                            <span> {{ bankStore.bank?.balance.toLocaleString("ru-RU") }} ₽</span>
                        </div>
                        <div class="bank-state up">
                            <span class="text-[16px]">
                                <ArrowUp />
                            </span>
                            <span>25% за 7д</span>
                        </div>
                    </div>
                </div>
            </div>
            </div>
        </div>
        <div class="info-block-2">
                <RouterLink :to="'/reports'" class="block">
                    <div class="block-header">
                        <div class="block-title">Отчеты</div>
                        <div class="text-[24px]">
                            <ArrowIcon />
                        </div>
                    </div>
                    <div class="block-info">
                        <div class="text-gray-500 text-base">Новый через 5д</div>
                        <div class="status">Отправлен</div>
                    </div>
                </RouterLink>
                <RouterLink :to="'/rules'" class="block">
                    <div class="block-header">
                        <div class="block-title">Правила</div>
                        <div class="text-[24px]">
                            <ArrowIcon />
                        </div>
                    </div>
                    <div class="block-info">
                        <div class="status">Смотреть</div>
                    </div>
                </RouterLink>
            </div>
    </section>
</template>

<style scoped>
.header {
    @apply flex flex-row justify-between pt-2 pr-4 pb-3 pl-4
}

.cur-user-name {
    @apply font-medium text-4xl tracking-[-1px]
}

.cur-user-benefits {
    @apply flex flex-row
}

.cur-user-pic {
    @apply h-[48px] w-[48px] rounded-[50%]
}

.notifications {
    @apply flex justify-center items-center h-[48px] w-[48px] rounded-[50%] bg-white text-white text-2xl -ml-2
}

.info-block-1 {
    @apply flex flex-row gap-1
}

.declaration-status {
    @apply flex flex-col gap-2 bg-white p-3 rounded-[16px]
    /* absolute w-[95%] bottom-[4px] left-1/2 -translate-x-1/2; */
}

.declaration-progress {
    font-family: "SF Pro Display";
    @apply text-[28px] leading-8
}

.progress-bar {
    @apply h-[24px] rounded-[20px] bg-[#ebebeb] overflow-hidden
}

.progress-fill {
    @apply h-[24px] bg-[#FF6644]
}

.members-and-bank {
    @apply flex flex-col justify-between gap-1
}


.info-block-2 {
    @apply flex flex-row gap-1
}

.block {
    @apply bg-white flex flex-col justify-between p-4 gap-6 rounded-[16px] min-w-[189px] min-h-[142px] flex-1
}

.block-header {
    font-family: "SF Pro Text";
    @apply flex flex-row justify-between font-medium
}

.block-info {
    font-family: "SF Pro Text";
    @apply flex flex-col tracking-[-0.4px] font-medium justify-center pr-5
}

.status {
    font-family: "SF Pro Display";
    @apply text-[28px] leading-8
}

.bank-state {
    @apply flex flex-row items-center gap-1 font-semibold tracking-[-0.25px] leading-5 text-sm
}

.up {
    @apply text-[#2EC055]
}

.inactive {
  opacity: 0.55;
}

</style>