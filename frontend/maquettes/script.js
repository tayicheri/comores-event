'use strict';
/*
 * Comores Events 2027 — démonstration du parcours billetterie.
 * Les prix sont des paramètres locaux modifiables. Aucun paiement n'est exécuté.
 * À brancher ultérieurement sur l'API de billetterie FOODYFLY.
 */
const EVENT_CONFIG = Object.freeze({
  name: 'Salon de la culture et de la gastronomie comorienne',
  year: 2027,
  dates: { samedi: 'Samedi 16 janvier', dimanche: 'Dimanche 17 janvier' },
  prices: { day: 20, pass: 30 }, // tarifs en euros par personne
  currency: 'EUR',
  maxTickets: 20
});

const state = { formula: 'day', day: 'samedi', quantity: 1 };
const formatEuro = value => new Intl.NumberFormat('fr-FR', { style: 'currency', currency: EVENT_CONFIG.currency, maximumFractionDigits: 0 }).format(value);
const $ = selector => document.querySelector(selector);

function updateTickets() {
  const unit = EVENT_CONFIG.prices[state.formula];
  const total = unit * state.quantity;
  document.querySelectorAll('.ticket-option').forEach(label => {
    label.classList.toggle('selected', label.dataset.formula === state.formula);
  });
  document.querySelectorAll('.day-choice').forEach(label => {
    label.classList.toggle('selected', label.querySelector('input').value === state.day);
  });
  const isDayTicket = state.formula === 'day';
  $('#day-picker').classList.toggle('disabled', !isDayTicket);
  $('#day-picker').disabled = !isDayTicket;
  $('#ticket-quantity').textContent = String(state.quantity);
  $('#ticket-total').textContent = formatEuro(total);
  $('#ticket-summary').textContent = `${state.quantity} billet${state.quantity > 1 ? 's' : ''} · ${formatEuro(unit)} / personne · ${isDayTicket ? EVENT_CONFIG.dates[state.day] : 'Pass samedi + dimanche'}`;
  $('#qty-minus').disabled = state.quantity <= 1;
  $('#qty-plus').disabled = state.quantity >= EVENT_CONFIG.maxTickets;
}

document.querySelectorAll('input[name="ticket-type"]').forEach(input => {
  input.addEventListener('change', () => {state.formula = input.value;updateTickets();});
});
document.querySelectorAll('input[name="ticket-day"]').forEach(input => {
  input.addEventListener('change', () => {state.day = input.value;updateTickets();});
});
$('#qty-minus').addEventListener('click', () => {state.quantity = Math.max(1, state.quantity - 1);updateTickets();});
$('#qty-plus').addEventListener('click', () => {state.quantity = Math.min(EVENT_CONFIG.maxTickets, state.quantity + 1);updateTickets();});

const modal = $('#demo-modal');
let lastFocused = null;
function closeModal(){modal.hidden = true;document.body.style.overflow = '';lastFocused?.focus();}
$('#checkout').addEventListener('click', () => {
  lastFocused = document.activeElement;
  const unit = EVENT_CONFIG.prices[state.formula];
  const formula = state.formula === 'day' ? `Billet 1 journée — ${EVENT_CONFIG.dates[state.day]}` : 'Pass 2 jours — samedi et dimanche';
  const summary = [formula, `${state.quantity} billet${state.quantity > 1 ? 's' : ''} × ${formatEuro(unit)}`, `Total : ${formatEuro(unit * state.quantity)}`];
  $('#modal-recap').replaceChildren(...summary.map(line => {const node=document.createElement('div');node.textContent=line;return node;}));
  modal.hidden = false;
  document.body.style.overflow = 'hidden';
  $('#modal-close').focus();
  // Production: au lieu d'afficher cette modale, appeler un endpoint serveur
  // créant une session de checkout dans l'API FOODYFLY puis rediriger.
});
$('#modal-close').addEventListener('click', closeModal);
$('#modal-done').addEventListener('click', closeModal);
modal.addEventListener('click', event => {if(event.target === modal) closeModal();});
document.addEventListener('keydown', event => {if(event.key === 'Escape' && !modal.hidden) closeModal();});

const mobileToggle = $('.mobile-toggle');
const navigation = $('#navigation');
mobileToggle.addEventListener('click', () => {
  const open = navigation.classList.toggle('open');
  mobileToggle.setAttribute('aria-expanded', String(open));
  mobileToggle.setAttribute('aria-label', open ? 'Fermer le menu' : 'Ouvrir le menu');
});
document.querySelectorAll('.main-nav a').forEach(link => link.addEventListener('click', () => {
  navigation.classList.remove('open');
  mobileToggle.setAttribute('aria-expanded', 'false');
  mobileToggle.setAttribute('aria-label', 'Ouvrir le menu');
}));
updateTickets();
