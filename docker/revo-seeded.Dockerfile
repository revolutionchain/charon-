FROM revo/revo

COPY ./fill_user_account.sh ./
COPY ./populate_and_run.sh ./
RUN chmod +x ./fill_user_account.sh
RUN chmod +x ./populate_and_run.sh

ENTRYPOINT [ "./populate_and_run.sh" ]