export type ProgramSession = {
  time: string
  title: string
  description?: string
  details?: string[]
  highlight?: boolean
}

export type ProgramDay = {
  key: 'samedi' | 'dimanche'
  label: string
  num: string
  month: string
  sessions: ProgramSession[]
}

export const PROGRAM_DAYS: readonly ProgramDay[] = Object.freeze([
  {
    key: 'samedi',
    label: 'SAMEDI',
    num: '16',
    month: 'JANVIER 2027',
    sessions: [
      {
        time: '10h30',
        title: 'Ouverture des stands',
        description: 'Accueil du public et découverte des espaces d’exposition.',
      },
      {
        time: '13h00',
        title: 'Cérémonie officielle d’ouverture',
        description: 'Allocutions des officiels des Comores et de l’UNESCO.',
        highlight: true,
      },
      {
        time: '14h00',
        title: 'Projection vidéo et exposition du patrimoine comorien',
        description:
          'Projection d’une vidéo sur l’inscription du bien comorien au patrimoine mondial de l’UNESCO, mettant en lumière les moments clés de sa validation en Corée du Sud.',
        details: [
          'Visite et exposition consacrées aux sites comoriens inscrits au patrimoine mondial de l’humanité, afin de faire découvrir leur histoire, leur richesse architecturale et leur valeur universelle exceptionnelle.',
        ],
      },
      {
        time: '14h30',
        title: 'Animations culturelles : traditions comoriennes',
        details: [
          'Mise en scène du Zifafa (mariage traditionnel comorien) : présentation des différentes étapes de la cérémonie, des coutumes, des tenues traditionnelles et des pratiques symboliques, à travers la représentation d’un couple de mariés.',
          'Danses traditionnelles comoriennes : Sambe et Chigoma.',
          'Danse traditionnelle anjouanaise (Tari) : présentation de cette danse du patrimoine culturel d’Anjouan, interprétée par des femmes originaires de Mutsamudu et de Domoni.',
        ],
      },
      {
        time: '16h00',
        title: 'Conférences-débats',
        description: 'Thème : « L’avenir du patrimoine culturel matériel et immatériel des Comores ».',
        details: [
          'Rencontre réunissant des spécialistes de la culture, de l’histoire et du patrimoine autour des enjeux de préservation, de transmission et de valorisation du patrimoine comorien.',
        ],
      },
      {
        time: '17h00',
        title: 'Ateliers culinaires et dégustation des anciens plats traditionnels comoriens',
        description:
          'Découverte et dégustation d’anciens plats traditionnels comoriens, mettant à l’honneur les recettes ancestrales et les savoir-faire culinaires transmis de génération en génération.',
      },
      {
        time: '22h00',
        title: 'Concert et animations culturelles',
        description:
          'Soirée festive autour de la musique comorienne, avec des prestations artistiques et des animations célébrant la richesse et la diversité du patrimoine musical de l’archipel.',
      },
    ],
  },
  {
    key: 'dimanche',
    label: 'DIMANCHE',
    num: '17',
    month: 'JANVIER 2027',
    sessions: [
      {
        time: '10h30',
        title: 'Ouverture des stands et expositions',
        description: 'Accueil du public, visite des stands et découverte des espaces d’exposition.',
      },
      {
        time: '12h00',
        title: 'Ateliers thématiques',
        description:
          'Restauration, préservation et valorisation touristique des sites comoriens inscrits au patrimoine mondial de l’UNESCO.',
        details: [
          'Les échanges porteront sur les enjeux de restauration et de préservation des sites, leur valorisation touristique, ainsi que les perspectives d’investissement et de mobilisation de la diaspora.',
          'Six associations représentant les villes abritant les sites inscrits au patrimoine mondial de l’UNESCO participeront aux échanges et aux réflexions sur l’avenir du patrimoine culturel des Comores.',
        ],
        highlight: true,
      },
      {
        time: '13h00',
        title: 'Animation culturelle : danse traditionnelle de Mohéli',
        description:
          'Présentation d’une danse traditionnelle de Mohéli par l’association ASDM (Association pour la Solidarité et le Développement de Monimoimdji), basée à Fomboni.',
      },
      {
        time: '14h00',
        title: 'Table ronde et atelier : artisanat et gastronomie comorienne',
        description:
          'Échanges autour de la richesse de l’artisanat et de la gastronomie comorienne, de leur transmission et de leur valorisation. L’atelier permettra de mettre en lumière les savoir-faire artisanaux et culinaires qui participent à la préservation et au rayonnement de la culture comorienne.',
      },
      {
        time: '15h00',
        title: 'Animation culturelle : danse traditionnelle féminine, le Wadaha',
        description:
          'Prestation mettant à l’honneur les rythmes, les gestes et les traditions culturelles comoriennes, interprétée par les femmes des associations HAWWA et AHDIR.',
      },
      {
        time: '22h00',
        title: 'Concert et cérémonie de clôture',
        description:
          'Grande soirée musicale pour célébrer la culture comorienne, avec des prestations artistiques et des animations mettant à l’honneur la richesse du patrimoine musical de l’archipel. Cette soirée marquera la clôture officielle de la deuxième édition du Salon de la culture et de la gastronomie comorienne en France.',
      },
    ],
  },
])
