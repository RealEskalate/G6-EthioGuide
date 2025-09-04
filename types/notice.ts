interface Procedure {
  id: string
  name: string
}

export default interface Notice {
  id: string
  orgId: string
  title: string
  body: string
  procedures: Procedure[]
  createdAt: string
  updatedAt: string
}