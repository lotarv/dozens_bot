<script lang="ts" setup>
import { useRoute } from 'vue-router'
import { useMembersStore } from '@/stores/membersStore'
import { Member } from '@/types/Member'
import ArrowLeft from '@/components/icons/ArrowLeft.vue'
import { declarations } from '@/mocks/declarations'
import DeclarationsComponent from '@/components/memberProfileView/DeclarationsComponent.vue'
const route = useRoute()
const membersStore = useMembersStore()
const username = route.params.username.toString()

const member: Member = membersStore.getMemberByUsername(username)
// const member: Member = {
//     fio: 'Андрей Грин',
//     username: "incetro",
//     annual_income: 100,
//     avatar_url: "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/e81977ef-8103-446d-9041-b5aa4fae82f5/alexander-hipp-iEEBWgY_6lA-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466T4P2WHZQ%2F20250501%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250501T114205Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjECMaCXVzLXdlc3QtMiJIMEYCIQDiXZSm8%2Ff0r%2ByAfcnhWdk1TtvZXBnaFPV7MVSViDpYwAIhALkYx1NSI6W1f91zsspoN5jQ0sRGtwdIosk0DImLYjd9KogECLz%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1IgyyUh18feFaj8cpzA0q3AM8wUcvNegErwOdN6HqaPN%2FQZYgwluFHCz2VDCcYSf%2Fdfc7hiv0NlXbJZPcVKVLXQhvGpQ1MMTuhxc0l5tqmWNmELhBaShR3UholvihwBrGT%2FreFgpQgAAlT2LegQkiIp1aR%2FhEpMDdYbg7u16cYYEMdBWI4zgVc8nMMgAMkeANeSHuaA9hU8ycYb9jY2tNGa7pdIpl30DKqO69IUxFP7zWtjpnW9pFwg8PftLmvX4XvDDN1S3%2B9Atii%2Bc1i%2FjexvLRh50WR4QXf2V8FnyPaYwyBaPaik5VDTHR5BAsounjTZ6ftOVp6Vv7gIzjqQrGsdf1WWfEPMAA%2BIoYVNqZ8MSOAdgZ3RsACk0WWCOM3%2BKusKCWy2JnhgCCzyU%2FIRPKgG2a68h9%2BV%2FMKugl1NC3jpDVJFWBI01fLNvoOKPfTW16XDlX5CniES59w44OMoY1Wovcy767zMTV4RwVDh%2BfBTV5an3Jckwyb2xCvPVJ5Ylu5FkEz6YLHG3LMYDBrbieshU2ZxPWwFM2DkhRdtuCqJEXI2enGSx9GdIZubn9MZcPpP%2FYSUd44XRQAWtjvRIOWBFzMPjP9GCBb2%2BUHC3sf%2B73n%2BHRcL63zPDnOgA7RJM4ff0kyAsIEbSl9iW%2FuTC4ns3ABjqkAb5NIxBEI1OvxtSSEbPAMmtbE167NVypqBrSueftmrHhHjIlCy3pNoDmktZsVtbIbfakUcGggaJVxtexXmrlq07Z9822cnLJP%2FSkQK7um92S1uJscaJA4%2BSIpcdQcsluxVs%2FLYPJjTF0BSJ10lIkFN3wa4C79vj1DG95D4hV2hcvC3sp1owGRtPpkzokXVB8eJ2QsCr%2Bd02TJY1nUNQAP%2FvZErsl&X-Amz-Signature=b585fc4eb747b9b2c6399a6ea56a931acde1ff8c75d2b17f959647473c01ae47&X-Amz-SignedHeaders=host&x-id=GetObject",
//     niche: "IT",
// }
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