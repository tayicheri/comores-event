export type Partner = {
  name: string
  src: string
}

export type PartnerGroup = {
  title: string
  partners: Partner[]
}

const LOGO_DIR = '/landing/partenaires'

export const PARTNER_GROUPS: readonly PartnerGroup[] = Object.freeze([
  {
    title: 'Partenaires institutionnels',
    partners: [
      { name: 'Union des Comores', src: `${LOGO_DIR}/union-des-comores.png` },
      { name: 'Délégation permanente de l’Union des Comores auprès de l’UNESCO', src: `${LOGO_DIR}/delegation-unesco-comores.jpg` },
      { name: 'UNESCO · Convention du patrimoine mondial', src: `${LOGO_DIR}/unesco-patrimoine-mondial.png` },
      { name: 'Office national du tourisme des Comores', src: `${LOGO_DIR}/office-tourisme-comores.png` },
      { name: 'Ville de Nanterre', src: '/landing/logo-ville-nanterre.svg' },
    ],
  },
  {
    title: 'Associations',
    partners: [
      { name: 'HAWWA', src: `${LOGO_DIR}/hawwa.png` },
      { name: 'AHDIR · Association humanitaire d’aide à l’intégration républicaine', src: `${LOGO_DIR}/ahdir.png` },
    ],
  },
  {
    title: 'Partenaires techniques',
    partners: [
      { name: 'Orange Money', src: `${LOGO_DIR}/orange-money.png` },
      { name: 'Yas', src: `${LOGO_DIR}/yas.png` },
      { name: 'Agence Presco', src: `${LOGO_DIR}/agence-presco.png` },
      { name: 'NYU', src: `${LOGO_DIR}/nyu.png` },
    ],
  },
])
