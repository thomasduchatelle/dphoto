export class User {
  constructor(public name: string,
              public email: string,
              public picture?: string) {
  }
}

export interface SecurityContextType {
  loggedUser?: User
}
