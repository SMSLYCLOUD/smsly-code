export interface User {
  id: number;
  username: string;
  email: string;
}

export interface Issue {
  id: number;
  repo_id: number;
  title: string;
  body: string;
  creator_id: number;
  creator: User;
  assignee_id?: number;
  assignee?: User;
  state: "open" | "closed";
  created_at: string;
  updated_at: string;
}

export interface Comment {
  id: number;
  issue_id: number;
  user_id: number;
  user: User;
  body: string;
  created_at: string;
}
