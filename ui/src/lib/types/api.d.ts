export namespace API {
  interface Base extends JSON {
    id: number;
  }

  export interface User extends Base {
    uid: number;
    name: string;
    display_name: string;
    avatar_url: string;
  }

  export interface Demo extends Base {
    source: string;
    status: string;
    map: string;
    players: DemoPlayer[];
    stats: StatsDemo;
    created_at: string;
    status_updated_at: string;
  }

  export interface DemoPlayer {
    user: User;
    stat: Stat;
    highlights?: Highlight[];
  }

  export interface Stat {
    result: string;
    start_team: string;
    kills: number;
    assists: number;
    deaths: number;
  }

  export interface StatsDemo {
    map: string;
    rounds_ct: number;
    rounds_t: number;
  }

  export interface Highlight extends Base {
    title: string;
  }
}
