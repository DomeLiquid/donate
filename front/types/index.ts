export interface GetProjectItem {
  project?: Project;
  user?: User;
}

export interface Project {
  pid: string;
  title: string;
  description: string;
  identityNumber: string;
  imgUrl: string;
  donateCnt: number;
  createdAt: string;
  link: string;
  user?: User;
}

export interface Asset {
  assetId: string;
  symbol: string;
  name: string;
  iconUrl: string;
  chainId: string;
  chainIconUrl: string;
  chainSymbol: string;
  priceUsd: string;
}

export interface User {
  identityNumber: string;
  fullName: string;
  mixinUid: string;
  avatarUrl: string;
  biography: string;
  mixinCreatedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface UserAction {
  identityNumber: string;
  fullName: string;
  avatarUrl: string;
  biography: string;
  assetId: string;
  amount: string;
  asset: Asset;
  project: Project;
  user?: User;
}
