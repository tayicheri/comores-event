export type HeritageSite = {
  /** `null` tant que la photo du site n’a pas été fournie : la carte affiche un emplacement neutre. */
  src: string | null
  alt: string
  name: string
  city: string
  island: string
  layout: string
}

const DIR = '/landing/patrimoine'

export const HERITAGE_SITES: readonly HeritageSite[] = Object.freeze([
  {
    src: `${DIR}/mosquee-vendredi-moroni.jpg`,
    alt: 'Ancienne mosquée du Vendredi de Moroni, blanche à arcades, au bord du port',
    name: 'Ancienne mosquée du Vendredi',
    city: 'Moroni',
    island: 'Grande Comore',
    layout: 'sm:col-span-2 lg:row-span-2',
  },
  {
    src: `${DIR}/porte-medina.jpg`,
    alt: 'Place Funihaziri à Iconi : porte blanche à arcades ouvrant sur une ruelle',
    name: 'Funihaziri',
    city: 'Iconi',
    island: 'Grande Comore',
    layout: 'row-span-2',
  },
  {
    src: `${DIR}/mosquee-domoni.jpg`,
    alt: 'Mihrab sculpté en pierre blanche à l’intérieur de la mosquée de Domoni',
    name: 'Mosquée de Domoni',
    city: 'Domoni',
    island: 'Anjouan',
    layout: 'row-span-2',
  },
  {
    src: `${DIR}/citadelle.jpg`,
    alt: 'Vue aérienne des remparts de la citadelle de Mutsamudu surplombant la ville et l’océan',
    name: 'Citadelle de Mutsamudu',
    city: 'Mutsamudu',
    island: 'Anjouan',
    layout: 'lg:col-span-2',
  },
  {
    src: `${DIR}/ruines-palais.jpg`,
    alt: 'Vestiges du palais Kaviridjewo en pierre volcanique avec son escalier',
    name: 'Palais Kaviridjewo',
    city: 'Iconi',
    island: 'Grande Comore',
    layout: '',
  },
  {
    src: null,
    alt: 'Arcades de pierre volcanique à Iconi (photo à venir)',
    name: 'Arcades de pierre volcanique',
    city: 'Iconi',
    island: 'Grande Comore',
    layout: '',
  },
  {
    src: null,
    alt: 'Portes et ruelles des médinas à Iconi (photo à venir)',
    name: 'Portes et ruelles des médinas',
    city: 'Iconi',
    island: 'Grande Comore',
    layout: '',
  },
  {
    src: `${DIR}/mnayfunu-itsandra.jpg`,
    alt: 'Façades anciennes en pierre volcanique du quartier de Mnayfunu à Itsandra',
    name: 'Mnayfunu',
    city: 'Itsandra',
    island: 'Grande Comore',
    layout: '',
  },
  {
    src: `${DIR}/ntsoudjini.jpg`,
    alt: 'Ruelle bordée de murs en pierre volcanique ouvrant sur l’océan à Ntsoudjini',
    name: 'Ruelles et remparts de Ntsoudjini',
    city: 'Ntsoudjini',
    island: 'Grande Comore',
    layout: '',
  },
  {
    src: `${DIR}/niches-palais-ujumbe-mutsamudu.jpg`,
    alt: 'Niches décorées sculptées dans le stuc blanc du palais Ujumbe',
    name: 'Niches décorées du palais Ujumbe',
    city: 'Mutsamudu',
    island: 'Anjouan',
    layout: '',
  },
])
