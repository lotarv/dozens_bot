<script lang="ts" setup>
import axios from 'axios';
import { ref, onBeforeMount, Ref } from 'vue';
import { getTelegramInitData } from '../services/auth';
import ArrowIcon from '@/components/icons/ArrowIcon.vue';
import {RouterLink} from 'vue-router'
interface Member {
    fio: string;
    avatar_url: string;
    niche: string;
    annual_income: number;
    username: string;
}

const members: Ref<Member[]> = ref([]);
// let members = [
//     {
//         "fio": "Анна",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/da1921d3-b856-4eb7-b607-febb8e03a713/lesha-tuman-jcZIhzgXuPI-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=c5c425afb50543c78b9f032ef5c8118540b45d725abda8103217911f54779409&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "Beauty",
//         "annual_income": 45,
//         "username": "anna"
//     },
//     {
//         "fio": "Дмитрий",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/2aa2aa2a-b845-40e4-b494-cd3816e13238/diego-hernandez-MSepzbKFz10-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=d5ab8af05dc49cd4a544d576da751b401cdc796445df1aa0799ba327916e656a&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "E-commerce",
//         "annual_income": 50,
//         "username": "dmitrii"
//     },
//     {
//         "fio": "Михаил Степанов",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/cf704917-3bd8-4942-9478-ad00a6896608/ivan-serediuk-1mEKyDTmTaw-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=f9a97abe6dc88ced5bd5698674936efdf679a6776712119fc5fea0166a62ad72&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "Education",
//         "annual_income": 55,
//         "username": "michail1"
//     },
//     {
//         "fio": "Александр Капустьянов",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/e5373e58-9dce-4bd8-b447-36f4d9f4e47c/alex-suprun-ZHvM3XIOHoE-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=c17c0e6a0b33c73e55bef125b1f64d42cd4ebdb40178e7935a3e1bf4427396e4&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "Food",
//         "annual_income": 60,
//         "username": "alexander"
//     },
//     {
//         "fio": "Анастасия Яновская",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/39475120-60bd-48e1-9013-d3ae6ea883ff/daniil-lobachev-97rfLsDp1RE-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=960f2ef2320275645089dce7ecafe925b3ce0fe0e700ad549823431df26af933&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "Home",
//         "annual_income": 65,
//         "username": "anastasia"
//     },
//     {
//         "fio": "Максим Борцов",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/8343ec76-c010-43e2-8eeb-205a0bfea0fd/midas-hofstra-a6PMA5JEmWE-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=80211a02b32dc34d3a62354cdbb9a0e6100ff803e223e6c060ccbb0a90c5f3a0&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "Clothes",
//         "annual_income": 70,
//         "username": "maksim"
//     },
//     {
//         "fio": "Юля Донцова",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/7ce152f3-ef7d-4587-8937-925c2c56bc97/seth-doyle-uJ8LNVCBjFQ-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=8789038ceb35680e9a4960da10ccf45d032105d02fd401d238de2751a5dc1cc4&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "Beauty",
//         "annual_income": 75,
//         "username": "julia"
//     },
//     {
//         "fio": "Полина HAPPYDAYS",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/09042a91-021e-4586-8bf3-aadc83c8574b/daniil-lobachev-XAo09LtQiAQ-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=9432a937151c8b190390116a4e515c736252127c35f53e1738c97e3a63a99670&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "Services",
//         "annual_income": 80,
//         "username": "polina"
//     },
//     {
//         "fio": "Юрий Терентьев",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/5d27ded4-8f22-4203-9394-9460b40bf7fd/clay-elliot-HfMCgqOLTyM-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=37e76c68630cfc9c813797beaa45bc2eaa4e60bd13ab7351746d04febbfbbc6e&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "FinTech",
//         "annual_income": 85,
//         "username": "yuri"
//     },
//     {
//         "fio": "Михаил Старостин",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/e81977ef-8103-446d-9041-b5aa4fae82f5/alexander-hipp-iEEBWgY_6lA-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=b34988861c3379ad2243a7fe730cf93c03d8d883ba365c7959bcafe3ebc236b4&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "E-commerce",
//         "annual_income": 90,
//         "username": "michail"
//     },
//     {
//         "fio": "Илья",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/0c493f69-7f75-4bab-8fb9-696cbb0a3e9c/midas-hofstra-tidSLv-UaNs-unsplash.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=f05de6cee340db0712e29c8165556b09f432040825b583a5eb5b6d6712a8e92b&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "Travel",
//         "annual_income": 95,
//         "username": "ilya"
//     },
//     {
//         "fio": "Андрей Грин",
//         "avatar_url": "https://prod-files-secure.s3.us-west-2.amazonaws.com/9a2e0635-b9d4-4178-a529-cf6b3bdce29d/dbd2a2a7-da84-44fa-b7f2-e87b8833c38f/photo_2025-04-29_10.09.55.jpeg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4665RIFGXAQ%2F20250429%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250429T162612Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEPj%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDFv%2FAzCBCUys%2Bgg6JqGMOCTaZK8QxVwKVMhWsS2moi0QIhAML%2BF5t5xBTsvuWJZ7gRKiKDuIPN%2FrDm8dg6CN3qUCYBKogECJH%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1Igx2Lyvejt9enn%2B9Rigq3AMRDNXulugxiL8%2F19ketSWE6cXl%2Bl1tyesTg8I5XTQaKzK8irP9d3L1QDIde1rrU6zEgfiFfJNHKrDdENo41MID0%2B%2BXYf68dtmo1P0zt9Bv5cOu%2BeibT2cDo1v%2Bvk1HBpzJd2hXnpPbOIyuBz4RmB3KdDXNIFHKdnOXiEgnDqE%2FfaFfAOQ9AuBVU1QPCzVZ47exVd6jdEGTGbJui939JkFYm6WbC8rFH0NMZ%2BnDThGZ%2Bdw%2FyQUAJrpI8MTADesk8jqna%2BgLdQZeO7pvjGcrt4vTjYqDjCk5eB%2FlKM2AtX0SHK6dDuzqSO3TlhSv77bvpeC7h4pB8It9jzcJKDDKPfpbnr7Aun9zd5CpPen%2F8ogDTojgQaP5tQH1LyylolzIvwfS%2Fl03xXeFr7Xm3gPAopjWL9JgYu%2F%2BlZkEz9CWokApLurQTx82Xy3DGsq1%2FR2XN42w63rIR2%2BdU%2B%2FwgPFwMiHRpp6dUGN0hri2P0O3eZZDjFqIvm%2F96hpyjB3i%2F0ajyvEU%2BXJHxOe1KKSfHupEog77g5CVyW5DkiUrGBTPTsViNzGeztGglXdQ6Wq3TmXH1ZuQS5yOYO61h2soPi4eOe9JM7tO9rQFv%2FqgNobtAErMPZQOuqxVO9PXtqQ90zCt98PABjqkAchi6fpn2wQiiGbKEdH9SfdDu%2F6q8lNsltAKEZlVXg5uQU3Ovgwb8Onlt7ZkHgt%2BmaMbmQsHn9%2FPhD34ldhjsOzRz%2FdwCaJnvH6y9vGr%2Fw9M4kwh94GyEb2r4QuAB%2FOEwKppSZmp038p%2BWoOfdmYVsnsMV3gBSBGUu0VEFlJoG0s0omSM1z7%2FPNB1kDq5bW2NDxfQGzlQatuLqdTkEM7AmaoAYwS&X-Amz-Signature=b4e6036dbcec1ad30e8876c20c481d73a486851e0422cc9bd98b1d11e4d53ccf&X-Amz-SignedHeaders=host&x-id=GetObject",
//         "niche": "IT",
//         "annual_income": 100,
//         "username": "incetro"
//     }
// ]
const isLoading = ref(false);
const error = ref<string | null>(null);

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

async function fetchMembers() {
    try {
        isLoading.value = true;
        error.value = null;

        const response = await axios.get<Member[]>(`${import.meta.env.VITE_API_URL}/members`, {
            headers: {
                'X-Telegram-Init-Data': getTelegramInitData(),
            },
        });
        members.value = response.data;
        console.log(members)
    } catch (err) {
        error.value = 'Failed to load members. Please try again later.';
        console.error('Failed to fetch members:', err);
    } finally {
        isLoading.value = false;
    }
}

onBeforeMount(async () => {
    await createOrUpdateUser();
    await fetchMembers();
});
</script>

<template>
    <section class="p-2 flex flex-col gap-4 font-medium">
        <h1 class="header">Участники</h1>
        <div v-if="isLoading" class="text-center text-xl">Загрузка...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="members.length > 0">
            <div class="members-container">
                <RouterLink
                :to="`/declarations/${member.username}`"
                class="member-card" v-for="member in members.reverse()" :key="member.username">
                    <div class="member-header">
                        <div class="flex items-center gap-2">
                            <img class="member-photo" :src="member.avatar_url" alt="">
                            <div class="member-declaration">
                                <span class="text-gray-500 tracking-[-0.4px]">Декларация</span>
                                <span>до 01.01.2001</span>
                            </div>
                        </div>
                        <div class="self-start text-2xl"><ArrowIcon/></div>
                    </div>
                    <div class="member-about">
                        <div class="member-name">{{ member.fio }}</div>
                        <div class="member-income">
                            <span>{{ member.annual_income }}М ₽ / год</span>
                            <div class="dot"></div>
                            <span class="niche">{{ member.niche }}</span>
                        </div>
                    </div>
                </RouterLink>
            </div>
        </div>
    </section>
</template>

<style scoped>
.header {
    @apply text-[36px] self-center font-medium tracking-[-1px]
}

.members-container{
    @apply flex flex-col gap-2 justify-center align-top
}

.member-card{
    @apply flex flex-col bg-white rounded-xl p-5 gap-4 cursor-pointer
}

.member-card{
    font-family: 'SF Pro Text', Roboto, emoji, sans-serif;
}

.member-header{
    @apply flex flex-row items-center justify-between w-full tracking-[-0.4px]
}

.member-photo{
    @apply rounded-[50%] w-[48px] h-[48px]
}

.member-declaration{
    @apply flex flex-col text-base
}

.members-about{
    @apply flex flex-col gap-1
}

.member-name {
    @apply text-xl tracking-[-0.4px]
}
.member-name {
    font-family: 'SF Pro Display', Roboto, emoji, sans-serif;
}
.member-income{
    @apply flex flex-row gap-1 text-[14px] font-[600] tracking-[-0.25px]
}
.dot {
    @apply w-[4px] h-[4px] bg-black inline-block rounded-[50%] self-center
}
</style>