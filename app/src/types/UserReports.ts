import { Report } from "./Report";
export interface UserReports {
    username: string;
    avatar_url: string;
    reports: Report[];
}