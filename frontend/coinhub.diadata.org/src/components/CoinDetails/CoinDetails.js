import axios from 'axios';
import numeral from 'numeral';
import moment from 'moment';
import router from '@/router';
import { AtomSpinner } from 'epic-spinners';
import sortBy from 'lodash/sortBy';
import { EventBus } from '@/main';


export default {
  components: {
    AtomSpinner
  },
  name: 'CoinDetails',
  props: {},
  data() {
    return {
      exchange_fields: [
        { key: 'Name', label: 'Exchange', sortable: true },
        { key: 'Price', label: 'Price', sortable: true },
        { key: 'Volume24', label: 'Volume (24h)', sortable: true },
        { key: 'Time', label: 'Last Updated', sortable: true },
        { key: 'show_trades', label: 'Trades', sortable: false },
      ],
      exchanges: [],
      last_trade_fields: [
        { key: 'Pair', label: 'Pair', sortable: true },
        { key: 'Price', label: 'Price', sortable: true },
        { key: 'Volume', label: 'Volume', sortable: true },
        { key: 'Time', label: 'Last Updated', sortable: true },
        { key: 'EstimatedUSDPrice', label: 'EstimatedUSDPrice', sortable: true },
      ],
      loading: true,
      errored: false,
      coinDetails:{},
      coinSymbol: '',
    };
  },
  created() {
    this.coinSymbol = this.$route.params.coinSymbol;
    EventBus.$emit('hideSearchInput', true);
  },
  mounted() {
    axios
    .get(`https://api.diadata.org/v1/symbol/${this.coinSymbol.toUpperCase()}`)
    .then((response) => {
      this.formatPairData(response.data);
    })
    .catch((error) => {
      console.log(error);
      this.errored = true;
    });
  },
  methods: {
  	formatPairData(data) {

      let {Coin, Change, Exchanges } = data;
      const change24 = (Coin.Price  - Coin.PriceYesterday) / Coin.PriceYesterday * 100;

      this.coinDetails = { 
          coinName: Coin.Name,
          coinSymbol: Coin.Symbol,
          coinPriceFormatted: Coin.Price < 1 ? '$'.concat(numeral(Coin.Price).format('0.0[0000]')) : '$'.concat(numeral(Coin.Price).format('0,0.00')),
          change24: change24,
          change24Formatted: change24 !== Number.POSITIVE_INFINITY ? numeral(change24).format('0,0.00').concat('%') : 'N/A',
          rank: this.$route.params.coinRank,
          volume24Formatted: '$'.concat(numeral(Coin.VolumeYesterdayUSD).format('0,0')),
          circulatingSupplyFormattedWithoutSymbol: numeral(Coin.CirculatingSupply).format('0,0'),
      };

      // format the exchanges
      Exchanges.forEach((exchange)=>{
        exchange.Price = '$'.concat(numeral(exchange.Price).format('0,0.00')),
        exchange.Volume24 = '$'.concat(numeral(exchange.VolumeYesterdayUSD).format('0,0')),
        exchange.Time = moment(exchange.Time).format("dddd, MMMM Do YYYY, h:mm:ss a");

      });

      Exchanges = sortBy(Exchanges, 'VolumeYesterdayUSD').reverse();

      this.exchanges = Exchanges;
      this.loading = false;
  	},
    switchCurrencies : function(currency){
      const { coinDataUSD, coinDataEUR } = coinData;
      if(currency === 'EUR'){
        this.coindata = coinDataEUR;
      }

      if(currency === 'USD'){
        this.coindata = coinDataUSD;
      }

      this.selectedCurrency = currency;

    }
  },
};