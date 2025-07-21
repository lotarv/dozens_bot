export interface TransactionMember {
  full_name: string;
  avatar_url: string;
}

export interface Transaction {
  member: TransactionMember;
  created_at: string;
  amount: number;
  reason: string;
}

export interface PiggyBank {
  balance: number;
  transactions: Transaction[];
}
