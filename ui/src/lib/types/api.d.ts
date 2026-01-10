export namespace API {
  interface Base extends JSON {
    id: number;
  }

  export interface User extends Base {
    uid: string;
    name: string;
    display_name: string;
    avatar_url: string;
  }

  export interface Demo extends Base {
    source: string;
    status: string;
    created_at: string;
    status_updated_at: string;
  }
}
